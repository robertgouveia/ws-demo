package http

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/robertgouveia/ws-demo/internal/server/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	addr         string
	echo         *echo.Echo
	errorChan    chan error
	shutdownChan chan struct{}
	ErrorLog     *log.Logger
	InfoLog      *log.Logger
	hub          *websocket.Hub
}

func (s Server) Shutdown() {
	close(s.shutdownChan)
}

func NewServer(addr string) *Server {
	return &Server{
		addr:         addr,
		errorChan:    make(chan error, 1),
		shutdownChan: make(chan struct{}),
		ErrorLog:     log.New(os.Stderr, fmt.Sprintf("[http %s]", addr), log.LstdFlags|log.Lshortfile),
		InfoLog:      log.New(os.Stdout, "[http] ", log.LstdFlags),
		echo:         echo.New(),
	}
}

func (s Server) Run() error {
	srv := &http.Server{
		Addr: s.addr,
	}

	gameManager := websocket.NewGameManager(s.InfoLog, s.ErrorLog)
	go gameManager.ListenForJoin()
	s.hub = websocket.NewHub(s.ErrorLog, s.InfoLog, gameManager)

	s.MapHandlers(s.echo, s.hub)

	go func() {
		if err := s.echo.StartServer(srv); err != nil {
			s.ErrorLog.Printf("Error Starting server: %s", err.Error())
			s.errorChan <- err
		}
	}()

	select {
	case e := <-s.errorChan:
		return e
	case <-s.shutdownChan:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.echo.Shutdown(ctx); err != nil {
			s.ErrorLog.Printf("Error Shutting down: %s", err.Error())
			return err
		}
		return nil
	}
}

func (s Server) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	s.InfoLog.Println("Shutting down server...")
	s.Shutdown()
}
