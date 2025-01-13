package application

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"giclo/internal/domain/models"
)

type Application struct {
	cfg *models.Config
}

func NewApplication(cfg *models.Config) *Application {
	app := &Application{cfg: cfg}
	return app
}

func getMails(servers []models.ServerConfig) (count uint32, tooltip []string, err error) {
	log.Debug().Msgf("Mail configuration is: `%v`", servers)
	var inboxCount uint32 = 0
	tooltipInfo := []string{}

	for _, srv := range servers {
		// TODO: parallel get for each server with error groups - v0.2?
		count, err := srv.MailCount()
		if err != nil {
			return 0, nil, err
		}
		inboxCount += count
		tooltipInfo = append(tooltipInfo, fmt.Sprintf("%s: %d", srv.Name, count))
	}
	return inboxCount, tooltipInfo, nil
}

func (app *Application) Start(ctx context.Context) {
	log.Debug().Msg("Starting mail application")
	inboxCount, tooltipInfo, err := getMails(app.cfg.Servers)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get mail count")
	}

	wo, err := models.NewWaybarOutput(inboxCount, tooltipInfo, app.cfg.EmptyInboxIcon, app.cfg.NonEmptyInboxIcon)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create Waybar output")
	}

	fmt.Println(wo)
	app.Stop(ctx)
}

func (app *Application) Stop(ctx context.Context) {
	log.Debug().Msg("Application stopped")
	os.Exit(0)
}
