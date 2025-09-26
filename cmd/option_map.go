package cmd

import (
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"
)

type OptionMap map[string]*interaction.CommandInteractionDataOption

func GenerateOptionMap(i *event.InteractionCreate) OptionMap {
	options := i.CommandData().Options
	optionMap := make(OptionMap, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func GenerateOptionMapForSubcommand(i *event.InteractionCreate) OptionMap {
	options := i.CommandData().Options[0].Options
	optionMap := make(OptionMap, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
