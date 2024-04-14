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
func (b *Bot) generalHandler(client *discordgo.Session, i *discordgo.InteractionCreate) {
	resp := utils.ResponseBuilder{
		I: i,
		C: client,
	}
	data := i.ApplicationCommandData()
	subInfo := data.Options[0]
	if subInfo == nil {
		err := resp.IsEphemeral().Message("No subcommand identified").Send()
		if err != nil {
			utils.SendAlert("general_commands_handler.go - No subcommand identified reply", err.Error())
			return
		}
		return
	}
	var cmd *Cmd
	for _, c := range b.Commands {
		if c.Name == data.Name {
			cmd = c
		}
	}
	if cmd == nil {
		err := resp.IsEphemeral().Message("Command not found").Send()
		if err != nil {
			utils.SendAlert("general_commands_handler.go - Command not found reply", err.Error())
			return
		}
	}
	if cmd.Subs == nil {
		utils.SendAlert("general_commands_handler.go - Checking subs", ErrSubsAreNil.Error())
		err := resp.IsEphemeral().Message("Internal error, please report it").Send()
		if err != nil {
			utils.SendAlert("general_commands_handler.go - Command not found reply", err.Error())
			return
		}
	}
	for _, s := range cmd.Subs {
		if subInfo.Name == s.Name {
			s.Handler(client, i)
		}
	}
	err := resp.IsEphemeral().Message("Subcommand not found").Send()
	if err != nil {
		utils.SendAlert("general_commands_handler.go - Subcommand not found reply", err.Error())
		return
	}
}
