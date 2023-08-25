package transport

import "github.com/nentenpizza/monopolium/wserver"

func (g *Game) OnStartTable(ctx *wserver.Context) error {
	client := ctx.Get("client").(*Client)
	if client == nil {
		return PlayerNotFoundErr
	}

	event := &EventJoinRoom{}
	err := ctx.Bind(event)
	if err != nil {
		return err
	}

	if client.Room() == nil {
		return RoomNotExistsErr
	}

	client.Room().Start()

	return nil
	//c := s.c.Read(client.Token.Username)
}
