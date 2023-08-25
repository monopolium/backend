package transport

import (
	"encoding/json"
	"sync"

	"github.com/nentenpizza/monopolium/wserver"

	"github.com/nentenpizza/monopolium/monopoly"

	log "github.com/sirupsen/logrus"
)

var Logger = log.New()

func init() {
	Logger.SetFormatter(&log.TextFormatter{})
}

type Map[T any] struct {
	entries map[string]*T
	sync.Mutex
}

func NewMap[T any]() *Map[T] {
	return &Map[T]{entries: map[string]*T{}}
}

func (m *Map[T]) Write(key string, value *T) {
	m.Lock()
	defer m.Unlock()
	m.entries[key] = value
}

func (m *Map[T]) Read(key string) *T {
	m.Lock()
	defer m.Unlock()
	return m.entries[key]
}

func (m *Map[T]) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.entries, key)
}

func (m *Map[T]) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(m.entries)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (m *Map[T]) ByID(key string) *T {
	m.Lock()
	defer m.Unlock()
	return m.entries[key]
}

type Game struct {
	Rooms   *Map[monopoly.Table]
	Clients *Map[Client]
	Secret  []byte
}

func NewGame(s []byte) *Game {
	return &Game{
		Rooms:   NewMap[monopoly.Table](),
		Clients: NewMap[Client](),
		Secret:  s,
	}
}

func (g *Game) Register(ws *wserver.Server) {
	ws.Use(g.WebsocketJWT)
	ws.Handle(EventTypeCreateRoom, g.OnCreateRoom)
	ws.Handle(EventTypeJoinRoom, g.OnJoinRoom)
	ws.Handle(EventTypeStartGame, g.OnStartTable)
	ws.Handle(wserver.OnOther, func(c *wserver.Context) error {
		log.Println(c.EventType(), c.Data())
		return nil
	})
}
