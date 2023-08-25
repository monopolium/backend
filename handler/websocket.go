package handler

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/nentenpizza/monopolium/app"

	"github.com/golang-jwt/jwt"

	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/monopolium/wserver"
)

type WebsocketHandlerGroup struct {
	Handler
	WServer *wserver.Server
	Secret  []byte
}

func (s WebsocketHandlerGroup) Register(h Handler, g *echo.Group) {
	s.Handler = h
	g.GET("/ws/:token", s.WebSocket)
	g.GET("/token", s.Token)
	s.WServer.Handle(wserver.OnConnect, s.OnConnect)
	s.WServer.Handle(wserver.OnDisconnect, s.OnDisconnect)
}

func (s WebsocketHandlerGroup) Token(c echo.Context) error {
	token := NewWithClaims(Claims{
		Username: strconv.Itoa(rand.Intn(10000)),
		ID:       int64(rand.Intn(10000)),
	})

	t, err := token.SignedString(s.Secret)
	if err != nil {
		return err
	}
	return c.JSON(200, echo.Map{"token": t})
}

func (s WebsocketHandlerGroup) WebSocket(c echo.Context) error {
	var token *jwt.Token
	var err error
	tok := c.Param("token")
	if tok == "" {
		return c.JSON(http.StatusBadRequest, app.Err("invalid token"))
	}
	token, err = jwt.ParseWithClaims(tok, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.Secret, nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad token")
	}
	if !token.Valid {
		return c.JSON(http.StatusBadRequest, "bad token")
	}
	ws, err := wserver.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("123")
		return err
	}
	s.WServer.Accept(ws, token)
	return nil
}

func (s WebsocketHandlerGroup) OnConnect(c *wserver.Context) error {
	log.Println("CONNECTED")
	return nil
}

func (s WebsocketHandlerGroup) OnDisconnect(c *wserver.Context) error {
	log.Println("DISCONNECTED")
	return nil
}
