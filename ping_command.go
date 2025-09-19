package gokord

import (
	"fmt"
	"time"

	cmd2 "github.com/anhgelus/gokord/cmd"
	discordgo "github.com/nyttikord/gokord"
)

func pingCommand(s *discordgo.Session, i *discordgo.InteractionCreate, _ cmd2.OptionMap, resp *cmd2.ResponseBuilder) {
	if err := resp.IsDeferred().Send(); err != nil { // sends the "is thinking..."
		s.LogError(err, "respond interaction")
		return
	}

	response, err := s.InteractionAPI().Response(i.Interaction)
	if err != nil {
		s.LogError(err, "interaction response")
	}

	var msg string

	timestamp, err := GetTimestampFromId(i.ID)
	if err != nil {
		s.LogError(err, "connect timestamp from ID")
		msg = ":ping_pong: Pong !"
	} else {
		s.LogDebug("%s", timestamp.Format(time.UnixDate))
		msg = fmt.Sprintf(
			":ping_pong: Pong !\nLatence du bot : `%d ms`\nLatence de l'API discord : `%d ms`",
			response.Timestamp.Sub(timestamp).Milliseconds(),
			s.HeartbeatLatency().Milliseconds(),
		)
	}

	if err = resp.SetMessage(msg).Send(); err != nil { // modifies the "is thinking..."
		s.LogError(err, "interaction response edit")
	}
}
