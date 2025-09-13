package gokord

import (
	"fmt"
	"time"

	cmd2 "github.com/anhgelus/gokord/cmd"
	"github.com/anhgelus/gokord/logger"
	discordgo "github.com/nyttikord/gokord"
)

func pingCommand(s *discordgo.Session, i *discordgo.InteractionCreate, _ cmd2.OptionMap, resp *cmd2.ResponseBuilder) {
	if err := resp.IsDeferred().Send(); err != nil { // sends the "is thinking..."
		logger.Alert("ping_command.go - Respond interaction", err.Error())
		return
	}

	response, err := s.InteractionAPI().Response(i.Interaction)
	if err != nil {
		logger.Alert("ping_command.go - Interaction response", err.Error())
	}

	var msg string

	timestamp, err := GetTimestampFromId(i.ID)
	if err != nil {
		logger.Alert("ping_command.go - Connect timestamp from ID", err.Error())
		msg = ":ping_pong: Pong !"
	} else {
		logger.Debug(timestamp.Format(time.UnixDate))
		msg = fmt.Sprintf(
			":ping_pong: Pong !\nLatence du bot : `%d ms`\nLatence de l'API discord : `%d ms`",
			response.Timestamp.Sub(timestamp).Milliseconds(),
			s.HeartbeatLatency().Milliseconds(),
		)
	}

	if err = resp.SetMessage(msg).Send(); err != nil { // modifies the "is thinking..."
		logger.Alert("ping_command.go - Interaction response edit", err.Error())
	}
}
