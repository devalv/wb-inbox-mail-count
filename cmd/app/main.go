package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"wb-inbox-mail-count/internal/app"
	"wb-inbox-mail-count/internal/config"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	log.Debug().Msgf("Config is: `%v`", cfg)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)
	defer cancel()

	app := app.NewApplication(cfg)
	go app.Start(ctx)
	<-ctx.Done()

	app.Stop(ctx)
}