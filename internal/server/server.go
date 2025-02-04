package server

import (
	root "CacheFlow/cmd"
	"CacheFlow/internal/database"
	"fmt"
	"net/http"
	"strconv"
)

type Server struct {
	port int
	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(root.Newflag.Port)
	NewServer := &Server{
		port: port,
		db: database.New(),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server
}
