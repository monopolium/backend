package monopoly

import "time"

const (
	EventTypeGameStarted        = "game_started"
	EventTypePlayerJoined       = "player_joined"
	EventTypePlayerLeft         = "player_left"
	EventTypeJoinAcknowledgment = "join_ack"

	EventTypeNextMove = "next_move"

	EventTypeTickerStarted = "ticker_started"

	EventTypeTradingStarted = "trading_started"
	EventTypeTradingEnded   = "trading_ended"

	EventTypePlayerLocationChanged = "player_location_changed"

	// EventTypePlayerReturnsToStart works with EventCash
	EventTypePlayerReturnsToStart = "player_made_circle"

	EventTypeWaitingForBuyDecision = "waiting_for_buy_decision"
	EventTypeAuctionStarted        = "auction_started"
	EventTypeAuctionEnded          = "auction_ended"
)

type Event struct {
	EventType string      `json:"event_type"`
	Data      interface{} `json:"data,omitempty"`
}

// EventPlayerJoined responsible for room joining acknowledgment of players
type EventPlayerJoined struct {
	Player *Player `json:"player"`
}

type EventPlayerNextMove struct {
	Player *Player `json:"player"`
}

// EventPlayerLeft responsible for room leaving acknowledgment of players
type EventPlayerLeft struct {
	PlayerID string `json:"player_id"`
}

// EventJoinAcknowledgment responsible for initializing base set of data for joined players
type EventJoinAcknowledgment struct {
	Players Players `json:"players"`
	Board   *Board  `json:"board"`
}

type EventTickerStarted struct {
	Duration time.Duration `json:"duration"`
}

type EventTradingStarted struct {
	Sender   *Player `json:"sender"`
	Receiver *Player `json:"receiver"`
}

type EventPlayerLocationChanged struct {
	PlayerID string `json:"playerID"`
	Location int    `json:"location"`
}

type EventCash struct {
	Cash int `json:"cash"`
}

type EventTradingEnded struct {
}

type EventGameStarted struct {
}
