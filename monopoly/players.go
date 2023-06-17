package monopoly

import "sync"

const(
    ColorRed = "red"
    ColorGreen = "green"
    ColorBlue = "blue"
    ColorYellow = "yellow"
)

type Player struct {
    sync.Mutex

    ID string  `json:"id"`
    Name string `json:"name"`

    // color is a color of an in-game figure of player representation
    color string

    updates chan *Event
}

func NewPlayer(ID string, name string) *Player {
    return &Player{ID: ID, Name: name, updates: make(chan *Event, 100)}
}

func (p *Player) Updates() <-chan *Event{
    return p.updates
}

func (p *Player) Color() string {
    return p.color
}