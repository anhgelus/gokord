package cmd

import (
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/interaction"
)

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate, optMap OptionMap, resp *ResponseBuilder)

type CommandBuilder interface {
	// SetHandler of the CommandBuilder (if it contains subcommand, it will never be called)
	SetHandler(handler CommandHandler) CommandBuilder
	// CanContainsSub makes the CommandBuilder able to contain subcommands
	CanContainsSub() CommandBuilder
	// AddSub to the CommandBuilder (also call ContainsSub)
	AddSub(s CommandBuilder) CommandBuilder
	// HasOption makes the CommandBuilder able to contain CommandOptionBuilder
	HasOption() CommandBuilder
	// AddOption to the CommandBuilder (also call HasOption)
	AddOption(s CommandOptionBuilder) CommandBuilder
	// AddContext to the command.
	// If it is empty, types.InteractionContextGuild will be added automatically
	AddContext(ctx types.InteractionContext) CommandBuilder
	// AddIntegrationType (where the command is installed).
	// If it is empty, types.IntegrationGuildInstall will be added automatically
	AddIntegrationType(ctx types.Integration) CommandBuilder
	// SetPermission of the CommandBuilder
	SetPermission(p *int64) CommandBuilder
	// GetName returns the name of the command
	GetName() string
	// HasSub returns true if the command has subcommands
	HasSub() bool
	// GetHandler returns the command's handler
	GetHandler() CommandHandler
	// GetSubs returns subcommands
	GetSubs() []CommandBuilder
	// ApplicationCommand returns the application command understandable by Discord
	ApplicationCommand() *interaction.Command
	setSub(bool)
	toSubCmd() *subCmd
}

type CommandOptionBuilder interface {
	// IsRequired informs that the CommandOptionBuilder is required
	IsRequired() CommandOptionBuilder
	// AddChoice to the CommandOptionBuilder
	AddChoice(ch CommandChoiceBuilder) CommandOptionBuilder
	toDiscordOption() *interaction.CommandOption
}

type CommandChoiceBuilder interface {
	toDiscordChoice() *interaction.CommandOptionChoice
}

// New creates a new CommandBuilder
func New(name string, description string) CommandBuilder {
	return &commandCreator{
		ContainsSub:      false,
		IsSub:            false,
		Name:             name,
		Contexts:         nil,
		IntegrationTypes: nil,
		Description:      description,
		Options:          []CommandOptionBuilder{},
		Subs:             []CommandBuilder{},
	}
}

// NewOption creates a new CommandOptionBuilder for CommandBuilder
func NewOption(t types.ApplicationCommandOption, name string, description string) CommandOptionBuilder {
	return &commandOptionCreator{
		Type:        t,
		Name:        name,
		Description: description,
		Required:    false,
		Choices:     []CommandChoiceBuilder{},
	}
}

// NewChoice creates a new CommandChoiceBuilder for CommandOptionBuilder
func NewChoice(name string, value interface{}) CommandChoiceBuilder {
	return &commandChoiceCreator{
		Name:  name,
		Value: value,
	}
}
