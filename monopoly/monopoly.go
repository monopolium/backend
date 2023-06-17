// Package monopoly contains all in-game logic
package monopoly

import (
	"math/rand"
	"sync"
)

var (
	colors []string = []string{ColorRed, ColorBlue, ColorGreen, ColorYellow}
)

type Players map[string]*Player

// Table represents a distinct game room with players and a Board to play
type Table struct {
	sync.Mutex
	players Players
	Board *Board
}

func (t *Table) broadcast(e *Event)  {
	t.Lock()
	defer t.Unlock()
	for _, p := range t.players{
		p.Lock()
		p.updates <- e
		p.Unlock()
	}
}

func (t *Table) AddPlayer(p *Player) {
	t.Lock()
	defer t.Unlock()
	color := colors[rand.Intn(len(colors)-1)]
	p.color = color
	t.players[p.ID] = p
}

// Board represents monopoly game board
type Board struct {

}