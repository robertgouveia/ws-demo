package http

import (
	"github.com/labstack/echo/v4"
	gHttp "github.com/robertgouveia/ws-demo/internal/game/delivery/http"
	"github.com/robertgouveia/ws-demo/internal/server/websocket"
)

func (s Server) MapHandlers(e *echo.Echo, h *websocket.Hub) {
	gameGroup := e.Group("/play")

	gameHandler := gHttp.NewGameHandler(h)
	gHttp.MapRoutes(gameGroup, gameHandler)
}
