package game

import "github.com/labstack/echo/v4"

type Handler interface {
	Connect() echo.HandlerFunc
}
