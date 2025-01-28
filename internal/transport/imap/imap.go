package imap

import (
	"strconv"

	"github.com/emersion/go-imap/v2"
	client "github.com/emersion/go-imap/v2/imapclient"
	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Name       string `yaml:"name"`
	Address    string `yaml:"address"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	UnreadOnly bool   `yaml:"unread_only"`
}

func (s ServerConfig) GetMailCount() (uint32, error) {
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
	if s.UnreadOnly {
		return unreadMailCount(c)
	}
	return mailCount(c)
}

func mailCount(c *client.Client) (uint32, error) {
	selectOptions := &imap.SelectOptions{ReadOnly: true}
	mbox, err := c.Select("INBOX", selectOptions).Wait()
	if err != nil {
		return 0, err
	}
	log.Debug().Msgf("INBOX contains %d messages", mbox.NumMessages)

	return mbox.NumMessages, nil
}

func unreadMailCount(c *client.Client) (uint32, error) {
	// без этого не прочитать SearchCriteria
	selectOptions := &imap.SelectOptions{ReadOnly: true}
	_, err := c.Select("INBOX", selectOptions).Wait()
	if err != nil {
		return 0, err
	}

	criteria := imap.SearchCriteria{NotFlag: []imap.Flag{imap.FlagSeen}}
	data, err := c.Search(&criteria, nil).Wait()
	if err != nil {
		return 0, err
	}

	if data.All.String() == "" {
		log.Debug().Msgf("INBOX contains 0 unread messages")
		return 0, nil
	}

	unreadCount, err := strconv.ParseUint(data.All.String(), 10, 32)
	if err != nil {
		return 0, err
	}
	log.Debug().Msgf("INBOX contains %d unread messages", uint32(unreadCount))
	return uint32(unreadCount), nil
}
