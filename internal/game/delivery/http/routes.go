package http

import (
	"github.com/labstack/echo/v4"
	"github.com/robertgouveia/ws-demo/internal/game"
)

func MapRoutes(e *echo.Group, h game.Handler) {
	e.GET("/ws", h.Connect())
}
