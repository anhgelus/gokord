package commands

import (
	"fmt"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
	"time"
)

func Ping(client *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := utils.ResponseBuilder{}
	err := resp.IsDeferred().Client(client).Interaction(i).Send()
	if err != nil {
		utils.SendAlert("ping.go - Respond interaction", err.Error())
	}
	resp.IsEdit()

	response, err := client.InteractionResponse(i.Interaction)
	if err != nil {
		utils.SendAlert("ping.go - Interaction response", err.Error())
	}

	var msg string

	interactionTimestamp, err := utils.GetTimestampFromId(i.ID)
	if err != nil {
		utils.SendAlert("ping.go - Connect timestamp from ID", err.Error())
		msg = ":ping_pong: Pong !"
	} else {
		utils.SendDebug(interactionTimestamp.Format(time.UnixDate))
		msg = fmt.Sprintf(
			":ping_pong: Pong !\nLatence du bot : `%d ms`\nLatence de l'API discord : `%d ms`",
			response.Timestamp.Sub(interactionTimestamp).Milliseconds(),
			client.HeartbeatLatency().Milliseconds(),
		)
	}
	err = resp.Message(msg).Send()

	if err != nil {
		utils.SendAlert("ping.go - Interaction response edit", err.Error())
	}
}
