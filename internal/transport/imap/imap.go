package imap

import (
	"strconv"

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

// TODO: отдельная функция для подключения, чтобы не дублировать код

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

func (s ServerConfig) UnreadMailCount() (uint32, error) {
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

	// без этого не прочитать SearchCriteria
	selectOptions := &imap.SelectOptions{ReadOnly: true}
	_, err = c.Select("INBOX", selectOptions).Wait()
	if err != nil {
		return 0, err
	}

	criteria := imap.SearchCriteria{NotFlag: []imap.Flag{imap.FlagSeen}}
	data, err := c.Search(&criteria, nil).Wait()
	if err != nil {
		return 0, err
	}

	if data.All.String() == "" {
		return 0, nil
	}

	unreadCount, err := strconv.ParseUint(data.All.String(), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(unreadCount), nil
}
