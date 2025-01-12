package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"

	"giclo/internal/domain/models"
)

func NewConfig(cfgPath string) (*models.Config, error) {
	var cfg models.Config

	err := cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("Config path is `%s`, Options are: `%v`", cfgPath, cfg)

	return &cfg, nil
}
