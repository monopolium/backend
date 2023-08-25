package transport

import (
	"errors"

	"github.com/nentenpizza/monopolium/handler"

	"github.com/golang-jwt/jwt"
	"github.com/nentenpizza/monopolium/wserver"
)

func (g *Game) WebsocketJWT(next wserver.HandlerFunc) wserver.HandlerFunc {
	return func(c *wserver.Context) error {
		tok := (c.Get("token")).(*jwt.Token)
		if tok == nil {
			return errors.New("token is nil")
		}
		token := handler.From(tok)
		client := g.Clients.Read(token.Username)
		if client != nil {
			client.conn = c.Conn
			if client.Token.Username == "" {
				client.Token = token
			}
		} else {
			client = NewClient(c.Conn, token, make([]interface{}, 0), make(chan bool))
			g.Clients.Write(client.Token.Username, client)
		}
		c.Set("client", client)
		return next(c)
	}
}
