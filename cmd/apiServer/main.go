package main

import (
	"myHttpServer/internal/http/rest/apiserver"
	"os"

	"log/slog"
)

func main() {

	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	server, err := apiserver.New(log)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if err := server.Start(); err != nil {
		log.Error(err.Error())
	}
}
