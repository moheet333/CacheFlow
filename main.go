/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"CacheFlow/cmd"
	"CacheFlow/internal/server"
	"fmt"
)

func main() {
	cmd.Execute()
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
