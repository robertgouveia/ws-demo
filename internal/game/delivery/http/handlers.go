package http

import (
	"github.com/labstack/echo/v4"
	"github.com/robertgouveia/ws-demo/internal/game"
	"github.com/robertgouveia/ws-demo/internal/server/websocket"
)

type GameHandler struct {
	hub *websocket.Hub
}

func NewGameHandler(hub *websocket.Hub) game.Handler {
	return &GameHandler{
		hub: hub,
	}
}

func (h GameHandler) Connect() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.hub.Handle(c.Response(), c.Request())
		return nil
	}
}
