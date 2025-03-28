package main

import (
	"context"
	"log"

	"github.com/vlone310/bss/config"
	db "github.com/vlone310/bss/internal/db/sqlc"
	"github.com/vlone310/bss/internal/http"
)

func main() {
	ctx := context.Background()
	config := config.MustLoadConfig(".")

	// setup persistance layer
	s := db.NewStore()
	s.Connect(ctx, config.DBSource)
	defer s.Close()

	srv := http.NewServer(s)
	log.Fatal(srv.ServeHTTP(config.ServerAddr))
}
