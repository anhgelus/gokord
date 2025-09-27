package gokord

import (
	"fmt"
	"time"

	cmd2 "github.com/anhgelus/gokord/cmd"
	"github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/event"
)

func pingCommand(s bot.Session, i *event.InteractionCreate, _ cmd2.OptionMap, resp *cmd2.ResponseBuilder) {
	if err := resp.IsDeferred().Send(); err != nil { // sends the "is thinking..."
		s.Logger().Error("respond interaction", "error", err)
		return
	}

	response, err := s.InteractionAPI().Response(i.Interaction)
	if err != nil {
		s.Logger().Error("interaction response", "error", err)
	}

	var msg string

	timestamp, err := GetTimestampFromId(i.ID)
	if err != nil {
		s.Logger().Error("connect timestamp from ID", "error", err)
		msg = ":ping_pong: Pong !"
	} else {
		s.Logger().Debug(timestamp.Format(time.UnixDate))
		msg = fmt.Sprintf(
			":ping_pong: Pong !\nLatence du bot : `%d ms`\nLatence de l'API discord : `%d ms`",
			response.Timestamp.Sub(timestamp).Milliseconds(),
			s.(*gokord.Session).HeartbeatLatency().Milliseconds(),
		)
	}

	if err = resp.SetMessage(msg).Send(); err != nil { // modifies the "is thinking..."
		s.Logger().Error("interaction response edit", "error", err)
	}
}
