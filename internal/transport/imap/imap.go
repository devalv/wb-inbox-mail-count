package imap

import (
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
