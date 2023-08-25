package handler

import (
	"github.com/labstack/echo/v4"
)

type EndpointGroup interface {
	Register(h Handler, g *echo.Group)
}

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Register(group *echo.Group, eg EndpointGroup) {
	eg.Register(h, group)
}
