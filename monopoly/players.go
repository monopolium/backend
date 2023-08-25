package monopoly

import "sync"

const (
	ColorRed    = "red"
	ColorGreen  = "green"
	ColorBlue   = "blue"
	ColorYellow = "yellow"
)

type Player struct {
	sync.Mutex

	ID   string `json:"id"`
	Name string `json:"name"`

	// color is a color of an in-game figure of player representation
	color string

	updates chan *Event

	// Index represents an order in game moves table
	Index int

	// Location represents index of tile that player currently standing on
	Location int `json:"location"`

	Property []Tile `json:"property"`
	Cash     int    `json:"cash"`
}

func NewPlayer(ID string, name string) *Player {
	return &Player{ID: ID, Name: name, updates: make(chan *Event, 100)}
}

func (p *Player) Updates() <-chan *Event {
	return p.updates
}

func (p *Player) Color() string {
	return p.color
}

func (p *Player) Send(e *Event) {
	p.Lock()
	defer p.Unlock()
	p.send(e)
}

func (p *Player) send(e *Event) {
	p.updates <- e
}
