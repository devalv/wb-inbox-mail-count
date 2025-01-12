package application

import (
	"context"
	"encoding/json"
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

type WaybarOutput struct {
	Text       string   `json:"text"`
	Tooltip    string   `json:"tooltip,omitempty"`
	Class      []string `json:"class,omitempty"`
	Percentage int      `json:"percentage"`
}

func NewApplication(cfg *models.Config) *Application {
	app := &Application{cfg: cfg}
	return app
}

func getMails(servers []models.ServerConfig) (*WaybarOutput, error) {
	// TODO: parallel get for each server with error groups
	log.Debug().Msgf("Mail configuration is: `%v`", servers)
	var inboxCount uint32 = 0

	for _, srvConfig := range servers {
		// TODO: another function with safe defer for logout
		c, err := client.DialTLS(srvConfig.Address, nil)
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("Connected to mail server `%s`", srvConfig.Address)
		if err := c.Login(srvConfig.Username, srvConfig.Password).Wait(); err != nil {
			return nil, err
		}

		log.Debug().Msgf("Logged in to mail server `%s`", srvConfig.Address)

		selectOptions := &imap.SelectOptions{ReadOnly: true}
		mbox, err := c.Select("INBOX", selectOptions).Wait()
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("INBOX contains %d messages", mbox.NumMessages)
		inboxCount += mbox.NumMessages
		c.Logout()
	}
	// TODO: в тултип выводить количество для каждого сервера
	wo := WaybarOutput{
		Text:       fmt.Sprintf("%d", inboxCount),
		Percentage: 100,
		// Class:      []string{"mail"},  // TODO: icon?
		Tooltip: fmt.Sprintf("%d", inboxCount),
	}
	return &wo, nil
}

func (app *Application) Start(ctx context.Context) {
	log.Debug().Msg("Starting mail application")
	wo, err := getMails(app.cfg.Servers)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get mail count") // TODO: err in domain.errors
	}
	log.Debug().Msgf("Waybar output is: `%v`", wo)

	// TODO: make another function
	str, _ := json.Marshal(wo)
	fmt.Println(string(str)) // output for waybar

	app.Stop(ctx)
}

func (app *Application) Stop(ctx context.Context) {
	log.Debug().Msg("Application stopped")
	os.Exit(0)
}
