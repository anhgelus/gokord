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
func (b *Bot) generalHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := utils.NewResponseBuilder(s, i)
	data := i.ApplicationCommandData()
	if len(data.Options) == 0 {
		utils.SendWarn("len(data.Options) == 0", "name", data.Name)
		err := resp.IsEphemeral().Message("No subcommand identified (may be a bug)").Send()
		if err != nil {
			utils.SendAlert("commands_handler.go - No subcommand identified reply", err.Error())
			return
		}
		return
	}
	subInfo := data.Options[0]
	if subInfo == nil {
		utils.SendWarn("subInfo == nil", "name", data.Name)
		err := resp.IsEphemeral().Message("No subcommand identified").Send()
		if err != nil {
			utils.SendAlert("commands_handler.go - No subcommand identified reply", err.Error())
			return
		}
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
		utils.SendWarn("cmd == nil", "name", data.Name)
		err := resp.IsEphemeral().Message("Command not found").Send()
		if err != nil {
			utils.SendAlert("commands_handler.go - Command not found reply", err.Error())
			return
		}
		return
	}
	if c.Subs == nil {
		utils.SendAlert("commands_handler.go - Checking subs", ErrSubsAreNil.Error())
		err := resp.IsEphemeral().Message("Internal error, please report it").Send()
		if err != nil {
			utils.SendAlert("commands_handler.go - Command not found reply", err.Error())
			return
		}
	}
	for _, sub := range c.Subs {
		if subInfo.Name == sub.Name {
			sub.Handler(s, i)
			return
		}
	}
	utils.SendWarn("Subcommand not found", "name", data.Name, "subInfo name", subInfo.Name)
	err := resp.IsEphemeral().Message("Subcommand not found").Send()
	if err != nil {
		utils.SendAlert("commands_handler.go - Subcommand not found reply", err.Error())
		return
	}
}
