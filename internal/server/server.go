package server

import (
	root "CacheFlow/cmd"
	"fmt"
	"net/http"
	"strconv"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(root.Newflag.Port)
	NewServer := &Server{
		port: port,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
