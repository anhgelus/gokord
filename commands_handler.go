package gokord

import (
	"errors"
	"github.com/anhgelus/gokord/cmd"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	ErrSubsAreNil = errors.New("subs are nil in general handler")
)

// generalHandler used for subcommand
func (b *Bot) generalHandler(s *discordgo.Session, i *discordgo.InteractionCreate, _ cmd.OptionMap, resp *cmd.ResponseBuilder) {
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
	var c cmd.CommandBuilder
	for _, cb := range b.Commands {
		if cb.GetName() == data.Name {
			c = cb
		}
	}
	if c == nil {
		sendWarn("cmd == nil", "Command not found")
		return
	}
	if c.GetSubs() == nil {
		utils.SendAlert("commands_handler.go - Checking subs", ErrSubsAreNil.Error())
		err := resp.IsEphemeral().SetMessage("Internal error, please report it").Send()
		if err != nil {
			utils.SendAlert("commands_handler.go - Internal error reply", err.Error())
		}
		return
	}
	for _, sub := range c.GetSubs() {
		if subInfo.Name == sub.GetName() {
			sub.GetHandler()(s, i, cmd.GenerateOptionMapForSubcommand(i), resp)
			return
		}
	}
	sendWarn("Subcommand not found", "Subcommand not found", "subInfo Name", subInfo.Name)
}
