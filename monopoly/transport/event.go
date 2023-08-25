package transport

const (
	EventTypeCreateRoom = "create_room"
	EventTypeJoinRoom   = "join_room"
	EventTypeLeaveRoom  = "leave_room"
	EventTypeStartGame  = "start_game"

	EventTypeThrowDice    = "throw_dice"
	EventTypeRequestTrade = "request_trade"
	EventTypeBuyTile      = "buy_tile"
)

type Event struct {
	EventType string      `json:"eventType"`
	Data      interface{} `json:"data"`
}

type (
	EventRoomCreated struct {
		ID string `json:"id"`
	}

	EventJoinRoom struct {
		ID string `json:"id"`
	}
)
