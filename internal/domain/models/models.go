package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/emersion/go-imap/v2"
	client "github.com/emersion/go-imap/v2/imapclient"
	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// TODO: куда поместить - adapters - v0.1?
func (s ServerConfig) MailCount() (uint32, error) {
	c, err := client.DialTLS(s.Address, nil)
	if err != nil {
		return 0, err
	}
	log.Debug().Msgf("Connected to mail server `%s`", s.Address)
	if err := c.Login(s.Username, s.Password).Wait(); err != nil {
		return 0, err
	}
	defer c.Logout()

	log.Debug().Msgf("Logged in to mail server `%s`", s.Address)

	selectOptions := &imap.SelectOptions{ReadOnly: true}
	mbox, err := c.Select("INBOX", selectOptions).Wait()
	if err != nil {
		return 0, err
	}
	log.Debug().Msgf("INBOX contains %d messages", mbox.NumMessages)

	return mbox.NumMessages, nil
}

type Config struct {
	Debug             bool           `yaml:"debug"`
	Servers           []ServerConfig `yaml:"servers"`
	EmptyInboxIcon    string         `yaml:"empty_inbox_icon"`
	NonEmptyInboxIcon string         `yaml:"non_empty_inbox_icon"`
}

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

// TODO: куда поместить - adapters - v0.1?
func (wo WaybarOutput) String() string {
	val, err := json.Marshal(wo)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal waybar output")
		return ""
	}
	return string(val)
}

// TODO: куда поместить - adapters - v0.1?
func NewWaybarOutput(inboxCount uint32, tooltipInfo []string, emptyInboxIcon, nonEmptyInboxIcon string) (WaybarOutput, error) {
	if inboxCount == 0 {
		return WaybarOutput{
			Text:    fmt.Sprintf("%d %s", inboxCount, emptyInboxIcon),
			Tooltip: strings.Join(tooltipInfo, "\n"),
		}, nil
	}

	return WaybarOutput{
		Text:    fmt.Sprintf("%d %s", inboxCount, nonEmptyInboxIcon),
		Tooltip: strings.Join(tooltipInfo, "\n"),
	}, nil
}
