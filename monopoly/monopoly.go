// Package monopoly contains all in-game logic
package monopoly

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	colors []string = []string{ColorRed, ColorBlue, ColorGreen, ColorYellow}
)

const (
	timeoutDuration = time.Second
	buyDuration     = time.Second
	auctionDuration = time.Second

	StateMoving  = "movingState"
	StateTrading = "tradingState"
	StateBuying  = "buyingState"
	StateAuction = "auctionState"
)

type Players map[string]*Player

// Table represents a distinct game room with players and a Board to play
type Table struct {
	sync.Mutex

	players          Players
	movingOrder      []*Player
	board            *Board
	currentMoveIndex int
	started          bool
	timeoutTicker    *time.Ticker
	cancelChan       chan struct{}

	state string
}

func NewTable(board *Board) *Table {
	return &Table{players: Players{}, board: board}
}

func (t *Table) Started() bool {
	return t.started
}

func (t *Table) Start() {
	t.Lock()
	defer t.Unlock()
	t.started = true
	for _, p := range t.players {
		p.Location = 0

		// defining player moving order, this is random because map iterations are random
		t.movingOrder = append(t.movingOrder, p)
	}

	// nextMove method will add 1 to this, so we need to make it -1 for the start of room
	t.currentMoveIndex = -1

	cancelChan := make(chan struct{}, 2)
	t.cancelChan = cancelChan
	t.timeoutTicker = t.enableTimeouts(cancelChan)

	t.nextMove()
}

func (t *Table) currentlyMovingPlayer() *Player {
	return t.movingOrder[t.currentMoveIndex]
}

func (t *Table) resetTicker(d time.Duration) {
	t.timeoutTicker.Reset(d)
	t.broadcast(&Event{EventType: EventTypeTickerStarted, Data: EventTickerStarted{Duration: d}})
}

func (t *Table) nextMove() {
	t.currentMoveIndex++

	if t.currentMoveIndex > len(t.movingOrder)-1 {
		t.currentMoveIndex = 0
	}

	t.broadcast(&Event{EventType: EventTypeNextMove, Data: EventPlayerNextMove{Player: t.currentlyMovingPlayer()}})
	t.state = StateMoving
	t.resetTicker(timeoutDuration)
}

func (t *Table) BuyTile(playerID string) error {
	t.Lock()
	defer t.Unlock()
	player := t.currentlyMovingPlayer()
	if playerID != player.ID {
		return errors.New("cannot move, not your turn")
	}

	tile := t.Located(player.Location)

	buyTile := func(tile ObtainableTile) {
		if player.Cash >= tile.Price() {
			player.Cash -= tile.Price()
			player.Property = append(player.Property, tile)
		}
	}

	if t, ok := tile.(ObtainableTile); ok {
		buyTile(t)
	}

	return nil
}

func (t *Table) Located(n int) Tile {
	return t.board.Body[n]
}

func (t *Table) Move(playerID string) error {
	t.Lock()
	defer t.Unlock()
	return t.move(playerID)
}

func (t *Table) move(playerID string) error {

	t.resetTicker(timeoutDuration)

	player := t.currentlyMovingPlayer()
	if playerID != player.ID {
		return errors.New("cannot move, not your turn")
	}

	firstDice, secondDice := t.rollDices()
	result := firstDice + secondDice

	newLocation := player.Location + result
	if newLocation > len(t.board.Body)-1 {
		newLocation = newLocation - len(t.board.Body) - 1
	}

	player.Location = newLocation

	t.broadcast(&Event{EventType: EventTypePlayerLocationChanged, Data: EventPlayerLocationChanged{
		PlayerID: player.ID,
		Location: player.Location,
	}})

	canSkip := t.applyLocationEffect(player)

	// we can go straight to the next move only if player can't do anything with tile
	if canSkip {
		t.nextMove()
	}
	return nil
}

func (t *Table) applyLocationEffect(player *Player) bool {
	if player.Location < 0 || player.Location > 39 {
		return true
	}
	tile := t.board.Body[player.Location]

	waitForBuyDecision := func() {
		t.state = StateBuying
		t.resetTicker(buyDuration)
		t.broadcast(&Event{EventType: EventTypeWaitingForBuyDecision})
	}

	switch tile.Type() {
	case TileTypeStart:
		player.Cash += 1000
		t.broadcast(&Event{
			EventType: EventTypePlayerReturnsToStart,
			Data:      EventCash{Cash: 1000},
		})
	case TileTypeCompany:
		waitForBuyDecision()
		return false
	case TileTypeAutomotive:
		waitForBuyDecision()
		return false
	case TileTypeService:
		waitForBuyDecision()
		return false
	}
	return true
}

func (t *Table) rollDices() (int, int) {
	return rand.Intn(6) + 1, rand.Intn(6) + 1
}

func (t *Table) enableTimeouts(cancel <-chan struct{}) *time.Ticker {
	ticker := time.NewTicker(timeoutDuration)
	go func() {
		for {
			select {
			case <-ticker.C:
				t.Timeout()
			case <-cancel:
				return
			}
		}
	}()
	ticker.Stop()
	return ticker
}

func (t *Table) Timeout() {
	t.Lock()
	defer t.Unlock()
	switch t.state {
	case StateMoving:
		_ = t.move(t.currentlyMovingPlayer().ID)
	case StateTrading:
		t.resetTicker(timeoutDuration)
		t.state = StateMoving
		t.broadcast(&Event{EventType: EventTypeTradingEnded, Data: EventTradingEnded{}})
	case StateBuying:
		t.resetTicker(auctionDuration)
		t.state = StateAuction
		t.broadcast(&Event{
			EventType: EventTypeAuctionStarted,
		})
	case StateAuction:
		t.nextMove()
	}
}

// Initialize initializes game board, boardBody is a json representation of the game board
func (t *Table) Initialize(boardBody []byte) {
	var board Board
	err := json.Unmarshal(boardBody, &board)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *Table) Broadcast(e *Event) {
	t.Lock()
	defer t.Unlock()
	t.broadcast(e)
}

func (t *Table) broadcast(e *Event) {
	for _, p := range t.players {
		p.Send(e)
	}
}

func (t *Table) AddPlayer(p *Player) {
	t.Lock()
	defer t.Unlock()
	if t.Started() {
		return
	}

	if _, ok := t.players[p.ID]; ok {
		return
	}

	color := colors[rand.Intn(len(colors)-1)]
	p.color = color
	t.players[p.ID] = p

	t.broadcast(&Event{
		EventType: EventTypePlayerJoined,
		Data:      EventPlayerJoined{Player: p},
	})

	p.send(&Event{
		EventType: EventTypeJoinAcknowledgment,
		Data: EventJoinAcknowledgment{
			Players: t.players,
			Board:   t.board,
		},
	})
}

func (t *Table) RemovePlayer(pID string) {
	t.Lock()
	defer t.Unlock()
	if t.Started() {
		return
	}

	delete(t.players, pID)

	t.broadcast(&Event{
		EventType: EventTypePlayerLeft,
		Data:      EventPlayerLeft{PlayerID: pID},
	})
}
