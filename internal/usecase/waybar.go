package usecase

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

func (wo WaybarOutput) String() string {
	val, err := json.Marshal(wo)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal waybar output")
		return ""
	}
	return string(val)
}

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
