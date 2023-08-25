package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/nentenpizza/monopolium/monopoly/transport"

	handler "github.com/nentenpizza/monopolium/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/monopolium/wserver"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.HideBanner = true

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Println(err)
		e.DefaultHTTPErrorHandler(err, c)
	}

	secret := []byte("d9799088-48bf-41c3-a109-6f09127f66bd")

	WServer := wserver.NewServer(wserver.Settings{
		Secret:  secret,
		OnError: func(e error, c *wserver.Context) { log.Println(e.Error(), c) },
		UseJWT:  true,
	})

	h := handler.NewHandler()

	h.Register(e.Group(""), handler.WebsocketHandlerGroup{WServer: WServer, Secret: secret})

	tr := transport.NewGame(secret)
	tr.Register(WServer)

	e.Logger.Fatal(e.Start(":7654"))
}
