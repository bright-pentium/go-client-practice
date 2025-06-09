package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bright-pentium/go-client-practice/internal/configs"
	server "github.com/bright-pentium/go-client-practice/internal/delivery/echo"
)

func main() {
	pgConfig, err := configs.LoadPostgresConfig("./configs/pg.env")
	if err != nil {
		panic(err)
	}

	appConfig, err := configs.LoadConfig("./configs/app.env")
	if err != nil {
		panic(err)
	}
	appConfig.DbURL = pgConfig.DBURL()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	server := server.NewServer(appConfig)
	if err := server.Serving(ctx); err != nil {
		panic(err)
	}
}
