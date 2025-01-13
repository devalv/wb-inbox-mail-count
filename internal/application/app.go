package application

import (
	"context"
	"fmt"
	"os"

	"github.com/emersion/go-imap/v2"
	client "github.com/emersion/go-imap/v2/imapclient"
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
	// TODO: parallel get for each server with error groups - v0.2?
	log.Debug().Msgf("Mail configuration is: `%v`", servers)
	var inboxCount uint32 = 0
	tooltipInfo := []string{}

	for _, srvConfig := range servers {
		// TODO: another function with safe defer for logout - v0.1
		c, err := client.DialTLS(srvConfig.Address, nil)
		if err != nil {
			return 0, nil, err
		}
		log.Debug().Msgf("Connected to mail server `%s`", srvConfig.Address)
		if err := c.Login(srvConfig.Username, srvConfig.Password).Wait(); err != nil {
			return 0, nil, err
		}

		log.Debug().Msgf("Logged in to mail server `%s`", srvConfig.Address)

		selectOptions := &imap.SelectOptions{ReadOnly: true}
		mbox, err := c.Select("INBOX", selectOptions).Wait()
		if err != nil {
			return 0, nil, err
		}
		log.Debug().Msgf("INBOX contains %d messages", mbox.NumMessages)
		inboxCount += mbox.NumMessages
		tooltipInfo = append(tooltipInfo, fmt.Sprintf("%s: %d", srvConfig.Name, mbox.NumMessages))
		c.Logout()
	}
	return inboxCount, tooltipInfo, nil
}

func (app *Application) Start(ctx context.Context) {
	log.Debug().Msg("Starting mail application")
	inboxCount, tooltipInfo, err := getMails(app.cfg.Servers)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get mail count")
	}

	wo, err := models.NewWaybarOutput(inboxCount, tooltipInfo)
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
