package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"giclo/internal/adapters/config"
	"giclo/internal/application"
	"giclo/internal/domain/models"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", path)
	}

	return nil
}

func parseFlags() (path string, err error) {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "./config.yml", "path to config file")
	flag.Parse()

	if err := validateConfigPath(cfgPath); err != nil {
		return "", err
	}

	return cfgPath, nil
}

func configureLogger(cfg *models.Config) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug mode enabled")
	}
}

func main() {
	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse flags")
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}
	configureLogger(cfg)

	log.Debug().Msgf("Config path is `%s`", cfgPath)
	log.Debug().Msgf("Config is: `%v`", cfg)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)
	defer cancel()

	app := application.NewApplication(cfg)
	go app.Start(ctx)
	<-ctx.Done()

	app.Stop(ctx)
}
