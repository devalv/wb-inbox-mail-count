package config

import (
	"github.com/ilyakaznacheev/cleanenv"

	"giclo/internal/domain/models"
)

const (
	EmptyInboxDefault    = "<span rise='2000'>󰶈</span>"
	NonEmptyInboxDefault = "<span color='#FF0000' rise='2000'>󰶍</span>"
)

func NewConfig(cfgPath string) (*models.Config, error) {
	var cfg models.Config

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.EmptyInboxIcon == "" {
		cfg.EmptyInboxIcon = EmptyInboxDefault
	}
	if cfg.NonEmptyInboxIcon == "" {
		cfg.NonEmptyInboxIcon = NonEmptyInboxDefault
	}
	return &cfg, nil
}
