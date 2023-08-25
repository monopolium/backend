package transport

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"

	"github.com/nentenpizza/monopolium/monopoly"
	"github.com/nentenpizza/monopolium/wserver"
)

func (g *Game) OnCreateRoom(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}
	if client.Room() != nil {
		return AlreadyInRoomErr
	}

	player := monopoly.NewPlayer(client.Token.Username, client.Token.Username)

	roomID := strconv.Itoa(rand.Intn(10000000000))

	var rawMap map[string]interface{}
	f, err := os.ReadFile("./monopoly/map.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(f, &rawMap)
	if err != nil {
		return err
	}

	b := &monopoly.Board{}

	b.Init(rawMap)

	table := monopoly.NewTable(b)

	g.Rooms.Write(roomID, table)

	table.AddPlayer(player)

	client.SetTable(table)
	client.SetChar(player)

	go client.ListenTable()

	err = client.WriteJSON(&Event{
		EventType: EventTypeCreateRoom,
		Data:      EventRoomCreated{roomID},
	})
	if err != nil {
		return err
	}

	err = client.WriteJSON(player)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) OnJoinRoom(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}

	event := &EventJoinRoom{}
	err := ctx.Bind(event)
	if err != nil {
		return err
	}
	//c := s.c.Read(client.Token.Username)

	if client.Room() != nil {
		return AlreadyInRoomErr
	}
	room := g.Rooms.Read(event.ID)
	if room == nil {
		return RoomNotExistsErr
	}
	if room.Started() {
		return RoomStartedErr
	}

	player := monopoly.NewPlayer(client.Token.Username, client.Token.Username)
	room.AddPlayer(player)

	client.SetTable(room)
	client.SetChar(player)

	go client.ListenTable()

	return client.WriteJSON(player)
}
