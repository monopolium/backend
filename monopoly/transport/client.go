package transport

import (
	"sync"
	"time"

	"github.com/nentenpizza/monopolium/handler"

	log "github.com/sirupsen/logrus"

	"github.com/nentenpizza/monopolium/monopoly"

	"github.com/nentenpizza/monopolium/wserver"
)

type Client struct {
	sync.Mutex
	conn *wserver.Conn
	*monopoly.Player
	table     *monopoly.Table
	AFK       bool
	Token     handler.Claims
	Unreached []interface{}
	quit      chan bool
	FloodWait time.Time
	EmojiWait time.Time
}

func (c *Client) LeaveRoom() {
	c.Lock()
	defer c.Unlock()
	c.table = nil
	c.Player = nil
}

func NewClient(conn *wserver.Conn, token handler.Claims, unreached []interface{}, quit chan bool) *Client {
	return &Client{
		conn:      conn,
		Token:     token,
		Unreached: unreached,
		quit:      quit,
		FloodWait: time.Now(),
		EmojiWait: time.Now(),
	}
}

func (c *Client) Conn() *wserver.Conn {
	c.Lock()
	defer c.Unlock()
	return c.conn
}

func (c *Client) UpdateConn(conn *wserver.Conn) {
	c.Lock()
	defer c.Unlock()
	c.conn = conn
}

func (c *Client) SetTable(r *monopoly.Table) {
	c.Lock()
	defer c.Unlock()
	c.table = r
}
func (c *Client) SetChar(plr *monopoly.Player) {
	c.Lock()
	defer c.Unlock()
	c.Player = plr
}

func (c *Client) Room() *monopoly.Table {
	c.Lock()
	defer c.Unlock()
	return c.table
}
func (c *Client) Char() *monopoly.Player {
	c.Lock()
	defer c.Unlock()
	return c.Player
}

func (c *Client) ListenTable() {
	for {
		if c.Player != nil {
			select {
			case value, ok := <-c.Player.Updates():
				if ok {
					c.WriteJSON(value)
				} else {
					return
				}
			case <-c.quit:
				return
			}
		} else {
			return
		}
	}
}

func (c *Client) WriteJSON(i interface{}) error {
	c.Lock()
	defer c.Unlock()
	var err error
	err = c.conn.WriteJSON(i)
	if err != nil {
		c.Unreached = append(c.Unreached, i)
		Logger.WithFields(log.Fields{
			"username":  c.Token.Username,
			"unreached": i,
			"room":      c.Room(),
		}).Info("Failed to send event")
	}
	return err
}

func (c *Client) SendUnreached() {
	if len(c.Unreached) > 0 {
		for _, e := range c.Unreached {
			c.WriteJSON(e)
			Logger.WithField(c.Token.Username, e).Info("Sent unreached")
		}
	}
	c.Unreached = make([]interface{}, 0)
}
