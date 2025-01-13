package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Debug   bool           `yaml:"debug"`
	Servers []ServerConfig `yaml:"servers"`
}

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

// TODO: Move to config fields - v0.2
const (
	EmptyInbox    = "<span rise='2000'>󰶈</span>"
	NonEmptyInbox = "<span color='#FF0000' rise='2000'>󰶍</span>"
)

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
func NewWaybarOutput(inboxCount uint32, tooltipInfo []string) (WaybarOutput, error) {
	if inboxCount == 0 {
		return WaybarOutput{
			Text:    fmt.Sprintf("%d %s", inboxCount, EmptyInbox),
			Tooltip: strings.Join(tooltipInfo, "\n"),
		}, nil
	}

	return WaybarOutput{
		Text:    fmt.Sprintf("%d %s", inboxCount, NonEmptyInbox),
		Tooltip: strings.Join(tooltipInfo, "\n"),
	}, nil
}
