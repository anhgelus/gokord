package cmd

import (
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
)

// subCmd is for the internal use of the API
type subCmd struct {
	*discordgo.ApplicationCommandOption
	Handler CommandHandler // Handler called
}

// commandCreator represents a generic command
type commandCreator struct {
	ContainsSub      bool
	IsSub            bool
	Name             string
	Permission       *int64
	Contexts         []discordgo.InteractionContextType
	IntegrationTypes []discordgo.ApplicationIntegrationType
	Description      string
	Options          []CommandOptionBuilder
	Subs             []CommandBuilder
	Handler          CommandHandler // Handler called
}

// commandOptionCreator represents a generic option of commandCreator
type commandOptionCreator struct {
	Type        discordgo.ApplicationCommandOptionType
	Name        string
	Description string
	Required    bool
	Choices     []CommandChoiceBuilder
}

// commandChoiceCreator represents a generic choice of commandOptionCreator
type commandChoiceCreator struct {
	Name  string
	Value interface{}
}

func (c *commandCreator) GetName() string {
	return c.Name
}

func (c *commandCreator) HasSub() bool {
	return c.ContainsSub
}

func (c *commandCreator) GetHandler() CommandHandler {
	return c.Handler
}

func (c *commandCreator) GetSubs() []CommandBuilder {
	return c.Subs
}

func (c *commandCreator) setSub(b bool) {
	c.IsSub = b
}

// SetHandler of the commandCreator (if commandCreator contains subcommand, it will never be called)
func (c *commandCreator) SetHandler(handler CommandHandler) CommandBuilder {
	c.Handler = handler
	return c
}

// CanContainsSub makes the commandCreator able to contain subcommands
func (c *commandCreator) CanContainsSub() CommandBuilder {
	c.ContainsSub = true
	c.Options = make([]CommandOptionBuilder, 0)
	return c
}

// AddSub to the commandCreator (also call ContainsSub)
func (c *commandCreator) AddSub(s CommandBuilder) CommandBuilder {
	c.CanContainsSub()
	s.setSub(true)
	c.Subs = append(c.Subs, s)
	return c
}

// HasOption makes the commandCreator able to contain commandOptionCreator
func (c *commandCreator) HasOption() CommandBuilder {
	c.ContainsSub = false
	c.Subs = make([]CommandBuilder, 0)
	return c
}

// AddOption to the commandCreator (also call HasOption)
func (c *commandCreator) AddOption(s CommandOptionBuilder) CommandBuilder {
	c.HasOption()
	c.Options = append(c.Options, s)
	return c
}

// AddContext to the command.
// If commandCreator.Contexts is empty, discordgo.InteractionContextGuild will be added automatically
func (c *commandCreator) AddContext(ctx discordgo.InteractionContextType) CommandBuilder {
	if c.Contexts == nil {
		c.Contexts = []discordgo.InteractionContextType{}
	}
	c.Contexts = append(c.Contexts, ctx)
	return c
}

func (c *commandCreator) AddIntegrationType(it discordgo.ApplicationIntegrationType) CommandBuilder {
	if c.IntegrationTypes == nil {
		c.IntegrationTypes = []discordgo.ApplicationIntegrationType{}
	}
	c.IntegrationTypes = append(c.IntegrationTypes, it)
	return c
}

// SetPermission of the commandCreator
func (c *commandCreator) SetPermission(p *int64) CommandBuilder {
	c.Permission = p
	return c
}

// Is returns true if the commandCreator is approximately the same as *discordgo.ApplicationCommand
func (c *commandCreator) Is(cmd *discordgo.ApplicationCommand) bool {
	return cmd.DefaultMemberPermissions == c.Permission &&
		cmd.Name == c.Name &&
		cmd.Description == c.Description &&
		len(cmd.Options) == len(c.Options)
}

// ApplicationCommand turns commandCreator into a *discordgo.ApplicationCommand
func (c *commandCreator) ApplicationCommand() *discordgo.ApplicationCommand {
	base := discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        c.Name,
		Description: c.Description,
	}
	if c.Permission != nil {
		base.DefaultMemberPermissions = c.Permission
	}
	if c.Contexts == nil || len(c.Contexts) == 0 {
		c.Contexts = []discordgo.InteractionContextType{discordgo.InteractionContextGuild}
	}
	base.Contexts = &c.Contexts
	if c.IntegrationTypes == nil || len(c.IntegrationTypes) == 0 {
		c.IntegrationTypes = []discordgo.ApplicationIntegrationType{discordgo.ApplicationIntegrationGuildInstall}
	}
	base.IntegrationTypes = &c.IntegrationTypes
	utils.SendDebug("Command creation", "name", c.Name, "has_sub", c.HasSub)
	if !c.ContainsSub {
		var options []*discordgo.ApplicationCommandOption
		for _, o := range c.Options {
			options = append(options, o.toDiscordOption())
		}
		base.Options = options
		return &base
	}
	var subs []*discordgo.ApplicationCommandOption
	for _, s := range c.Subs {
		sub := s.toSubCmd()
		subs = append(subs, sub.ApplicationCommandOption)
	}
	base.Options = subs
	return &base
}

// ToSubCmd turns commandCreator into a subCmd
func (c *commandCreator) toSubCmd() *subCmd {
	base := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        c.Name,
		Description: c.Description,
	}
	utils.SendDebug("Subcommand creation", "name", c.Name, "len(options)", len(c.Options))
	if len(c.Options) > 0 {
		var options []*discordgo.ApplicationCommandOption
		for _, o := range c.Options {
			options = append(options, o.toDiscordOption())
		}
		base.Options = options
	}
	return &subCmd{
		ApplicationCommandOption: &base,
		Handler:                  c.Handler,
	}
}

// IsRequired informs that the commandOptionCreator is required
func (o *commandOptionCreator) IsRequired() CommandOptionBuilder {
	o.Required = true
	return o
}

// AddChoice to the commandOptionCreator
func (o *commandOptionCreator) AddChoice(c CommandChoiceBuilder) CommandOptionBuilder {
	o.Required = true
	o.Choices = append(o.Choices, c)
	return o
}

// toDiscordOption turns commandOptionCreator into a discordgo.ApplicationCommandOption
func (o *commandOptionCreator) toDiscordOption() *discordgo.ApplicationCommandOption {
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, c := range o.Choices {
		choices = append(choices, c.toDiscordChoice())
	}
	return &discordgo.ApplicationCommandOption{
		Type:        o.Type,
		Name:        o.Name,
		Description: o.Description,
		Required:    o.Required,
		Choices:     choices,
	}
}

// toDiscordChoice turns commandChoiceCreator into a discordgo.ApplicationCommandOptionChoice (internal use of the API only)
func (c *commandChoiceCreator) toDiscordChoice() *discordgo.ApplicationCommandOptionChoice {
	return &discordgo.ApplicationCommandOptionChoice{
		Name:  c.Name,
		Value: c.Value,
	}
}
