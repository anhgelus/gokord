package commands

import (
	"fmt"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
	"time"
)

func Ping(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := utils.NewResponseBuilder(s, i)
	if err := resp.IsDeferred().Send(); err != nil { // sends the "is thinking..."
		utils.SendAlert("ping.go - Respond interaction", err.Error())
	}

	response, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		utils.SendAlert("ping.go - Interaction response", err.Error())
	}

	var msg string

	timestamp, err := utils.GetTimestampFromId(i.ID)
	if err != nil {
		utils.SendAlert("ping.go - Connect timestamp from ID", err.Error())
		msg = ":ping_pong: Pong !"
	} else {
		utils.SendDebug(timestamp.Format(time.UnixDate))
		msg = fmt.Sprintf(
			":ping_pong: Pong !\nLatence du bot : `%d ms`\nLatence de l'API discord : `%d ms`",
			response.Timestamp.Sub(timestamp).Milliseconds(),
			s.HeartbeatLatency().Milliseconds(),
		)
	}

	if err = resp.SetMessage(msg).Send(); err != nil { // modifies the "is thinking..."
		utils.SendAlert("ping.go - Interaction response edit", err.Error())
	}
}
