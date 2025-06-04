package gokord

import "github.com/bwmarrin/discordgo"

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type CommandBuilder interface {
	// SetHandler of the CommandBuilder (if it contains subcommand, it will never be called)
	SetHandler(handler CommandHandler) CommandBuilder
	// ContainsSub makes the CommandBuilder able to contain subcommands
	ContainsSub() CommandBuilder
	// AddSub to the CommandBuilder (also call ContainsSub)
	AddSub(s CommandBuilder) CommandBuilder
	// HasOption makes the CommandBuilder able to contain CommandOptionBuilder
	HasOption() CommandBuilder
	// AddOption to the CommandBuilder (also call HasOption)
	AddOption(s CommandOptionBuilder) CommandBuilder
	// AddContext to the command.
	// If it is empty, discordgo.InteractionContextGuild will be added automatically
	AddContext(ctx discordgo.InteractionContextType) CommandBuilder
	// AddIntegrationType (where the command is installed).
	// If it is empty, discordgo.ApplicationIntegrationGuildInstall will be added automatically
	AddIntegrationType(ctx discordgo.ApplicationIntegrationType) CommandBuilder
	// SetPermission of the CommandBuilder
	SetPermission(p *int64) CommandBuilder
	toCreator() *commandCreator
}

type CommandOptionBuilder interface {
	// IsRequired informs that the CommandOptionBuilder is required
	IsRequired() CommandOptionBuilder
	// AddChoice to the CommandOptionBuilder
	AddChoice(ch CommandChoiceBuilder) CommandOptionBuilder
	toCreator() *commandOptionCreator
}

type CommandChoiceBuilder interface {
	toCreator() *commandChoiceCreator
}

type commandBuilderCreator struct {
	*commandCreator
}

type commandOptionBuilderCreator struct {
	*commandOptionCreator
}

type commandChoiceBuilderCreator struct {
	*commandChoiceCreator
}

func (c *commandBuilderCreator) SetHandler(handler CommandHandler) CommandBuilder {
	c.commandCreator.SetHandler(handler)
	return c
}

func (c *commandBuilderCreator) ContainsSub() CommandBuilder {
	c.commandCreator.ContainsSub()
	return c
}

func (c *commandBuilderCreator) AddSub(s CommandBuilder) CommandBuilder {
	c.commandCreator.AddSub(s.toCreator())
	return c
}

func (c *commandBuilderCreator) HasOption() CommandBuilder {
	c.commandCreator.HasOption()
	return c
}

func (c *commandBuilderCreator) AddOption(s CommandOptionBuilder) CommandBuilder {
	c.commandCreator.AddOption(s.toCreator())
	return c
}

func (c *commandBuilderCreator) AddContext(ctx discordgo.InteractionContextType) CommandBuilder {
	c.commandCreator.AddContext(ctx)
	return c
}

func (c *commandBuilderCreator) AddIntegrationType(ctx discordgo.ApplicationIntegrationType) CommandBuilder {
	c.commandCreator.AddIntegrationType(ctx)
	return c
}

func (c *commandBuilderCreator) SetPermission(p *int64) CommandBuilder {
	c.commandCreator.SetPermission(p)
	return c
}

func (c *commandBuilderCreator) toCreator() *commandCreator {
	return c.commandCreator
}

func (c *commandOptionBuilderCreator) IsRequired() CommandOptionBuilder {
	c.commandOptionCreator.IsRequired()
	return c
}

func (c *commandOptionBuilderCreator) AddChoice(ch CommandChoiceBuilder) CommandOptionBuilder {
	c.commandOptionCreator.AddChoice(ch.toCreator())
	return c
}

func (c *commandOptionBuilderCreator) toCreator() *commandOptionCreator {
	return c.commandOptionCreator
}

func (c *commandChoiceBuilderCreator) toCreator() *commandChoiceCreator {
	return c.commandChoiceCreator
}

// NewCommand creates a new CommandBuilder
func NewCommand(name string, description string) CommandBuilder {
	return &commandBuilderCreator{
		commandCreator: &commandCreator{
			HasSub:           false,
			IsSub:            false,
			Name:             name,
			Contexts:         nil,
			IntegrationTypes: nil,
			Description:      description,
			Options:          []*commandOptionCreator{},
			Subs:             []*commandCreator{},
		},
	}
}

// NewOption creates a new CommandOptionBuilder for CommandBuilder
func NewOption(t discordgo.ApplicationCommandOptionType, name string, description string) CommandOptionBuilder {
	return &commandOptionBuilderCreator{
		commandOptionCreator: &commandOptionCreator{
			Type:        t,
			Name:        name,
			Description: description,
			Required:    false,
			Choices:     []*commandChoiceCreator{},
		},
	}
}

// NewChoice creates a new CommandChoiceBuilder for CommandOptionBuilder
func NewChoice(name string, value interface{}) CommandChoiceBuilder {
	return &commandChoiceBuilderCreator{
		commandChoiceCreator: &commandChoiceCreator{
			Name:  name,
			Value: value,
		},
	}
}
