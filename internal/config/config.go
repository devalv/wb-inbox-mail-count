package config

import (
	"flag"
	"fmt"
	"os"

	"wb-inbox-mail-count/internal/domain/consts"
	"wb-inbox-mail-count/internal/transport/imap"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Debug             bool                `yaml:"debug"`
	Servers           []imap.ServerConfig `yaml:"servers"`
	EmptyInboxIcon    string              `yaml:"empty_inbox_icon"`
	NonEmptyInboxIcon string              `yaml:"non_empty_inbox_icon"`
	ConfigPath        string
}

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

func NewConfig() (*Config, error) {
	cfgPath, err := parseFlags()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse flags")
	}

	var cfg Config
	err = cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		return nil, err
	}
	cfg.ConfigPath = cfgPath
	cfg.ConfigureLogger()
	if cfg.EmptyInboxIcon == "" {
		cfg.EmptyInboxIcon = consts.EmptyInboxDefault
	}
	if cfg.NonEmptyInboxIcon == "" {
		cfg.NonEmptyInboxIcon = consts.NonEmptyInboxDefault
	}
	return &cfg, nil
}

func (cfg *Config) ConfigureLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug mode enabled")
	}
}
