package gokord

import (
	"errors"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	ErrSubsAreNil = errors.New("subs are nil in general handler")
)

// generalHandler used for subcommand
func (b *Bot) generalHandler(s *discordgo.Session, i *discordgo.InteractionCreate, optMap utils.OptionMap, resp *utils.ResponseBuilder) {
	data := i.ApplicationCommandData()
	sendWarn := func(msg string, msgSend string, more ...interface{}) {
		utils.SendWarn(msg, "name", data.Name, more)
		err := resp.IsEphemeral().SetMessage(msgSend).Send()
		if err != nil {
			utils.SendAlert("commands_handler.go - "+msgSend+" reply", err.Error())
		}
	}
	if len(data.Options) == 0 {
		sendWarn("len(data.Options) == 0", "No subcommand identified (may be a bug)")
		return
	}
	subInfo := data.Options[0]
	if subInfo == nil {
		sendWarn("subInfo == nil", "No subcommand identified")
		return
	}
	var c *cmd
	for _, cb := range b.Commands {
		cr := cb.toCreator()
		if cr.Name == data.Name {
			c = cr.ToCmd()
		}
	}
	if c == nil {
		sendWarn("cmd == nil", "Command not found")
		return
	}
	if c.Subs == nil {
		utils.SendAlert("commands_handler.go - Checking subs", ErrSubsAreNil.Error())
		err := resp.IsEphemeral().SetMessage("Internal error, please report it").Send()
		if err != nil {
			utils.SendAlert("commands_handler.go - Internal error reply", err.Error())
		}
		return
	}
	for _, sub := range c.Subs {
		if subInfo.Name == sub.Name {
			sub.Handler(s, i, utils.GenerateOptionMapForSubcommand(i), resp)
			return
		}
	}
	sendWarn("Subcommand not found", "Subcommand not found", "subInfo Name", subInfo.Name)
}
