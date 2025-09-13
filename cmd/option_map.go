package cmd

import (
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/interaction"
)

type OptionMap map[string]*interaction.CommandInteractionDataOption

func GenerateOptionMap(i *discordgo.InteractionCreate) OptionMap {
	options := i.CommandData().Options
	optionMap := make(OptionMap, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func GenerateOptionMapForSubcommand(i *discordgo.InteractionCreate) OptionMap {
	options := i.CommandData().Options[0].Options
	optionMap := make(OptionMap, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
