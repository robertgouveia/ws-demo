package main

import (
	server "github.com/robertgouveia/ws-demo/internal/server/http"
)

func main() {
	s := server.NewServer(":4000")

	s.InfoLog.Printf("Starting server on port: %s", ":4000")
	go s.ListenForShutdown()
	if err := s.Run(); err != nil {
		panic(err)
	}
}
