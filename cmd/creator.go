package cmd

import (
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/interaction"
	"github.com/nyttikord/gokord/logger"
)

// subCmd is for the internal use of the API
type subCmd struct {
	*interaction.CommandOption
	Handler CommandHandler // Handler called
}

// commandCreator represents a generic command
type commandCreator struct {
	logger.Logger
	ContainsSub      bool
	IsSub            bool
	Name             string
	Permission       *int64
	Contexts         []types.InteractionContext
	IntegrationTypes []types.IntegrationInstall
	Description      string
	Options          []CommandOptionBuilder
	Subs             []CommandBuilder
	Handler          CommandHandler // Handler called
}

// commandOptionCreator represents a generic option of commandCreator
type commandOptionCreator struct {
	Type        types.CommandOption
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
// If commandCreator.Contexts is empty, types.InteractionContextGuild will be added automatically
func (c *commandCreator) AddContext(ctx types.InteractionContext) CommandBuilder {
	if c.Contexts == nil {
		c.Contexts = []types.InteractionContext{}
	}
	c.Contexts = append(c.Contexts, ctx)
	return c
}

func (c *commandCreator) AddIntegrationType(it types.IntegrationInstall) CommandBuilder {
	if c.IntegrationTypes == nil {
		c.IntegrationTypes = []types.IntegrationInstall{}
	}
	c.IntegrationTypes = append(c.IntegrationTypes, it)
	return c
}

// SetPermission of the commandCreator
func (c *commandCreator) SetPermission(p *int64) CommandBuilder {
	c.Permission = p
	return c
}

// Is returns true if the commandCreator is approximately the same as *interaction.Command
func (c *commandCreator) Is(cmd *interaction.Command) bool {
	return cmd.DefaultMemberPermissions == c.Permission &&
		cmd.Name == c.Name &&
		cmd.Description == c.Description &&
		len(cmd.Options) == len(c.Options)
}

// ApplicationCommand turns commandCreator into a *interaction.Command
func (c *commandCreator) ApplicationCommand() *interaction.Command {
	base := interaction.Command{
		Type:        types.CommandChat,
		Name:        c.Name,
		Description: c.Description,
	}
	if c.Permission != nil {
		base.DefaultMemberPermissions = c.Permission
	}
	if c.Contexts == nil || len(c.Contexts) == 0 {
		c.Contexts = []types.InteractionContext{types.InteractionContextGuild}
	}
	base.Contexts = &c.Contexts
	if c.IntegrationTypes == nil || len(c.IntegrationTypes) == 0 {
		c.IntegrationTypes = []types.IntegrationInstall{types.IntegrationInstallGuild}
	}
	base.IntegrationTypes = &c.IntegrationTypes
	c.LogDebug("Command creation", "name", c.Name, "has_sub", c.HasSub)
	if !c.ContainsSub {
		var options []*interaction.CommandOption
		for _, o := range c.Options {
			options = append(options, o.toDiscordOption())
		}
		base.Options = options
		return &base
	}
	var subs []*interaction.CommandOption
	for _, s := range c.Subs {
		sub := s.toSubCmd()
		subs = append(subs, sub.CommandOption)
	}
	base.Options = subs
	return &base
}

// ToSubCmd turns commandCreator into a subCmd
func (c *commandCreator) toSubCmd() *subCmd {
	base := interaction.CommandOption{
		Type:        types.CommandOptionSubCommand,
		Name:        c.Name,
		Description: c.Description,
	}
	c.LogDebug("Subcommand creation", "name", c.Name, "len(options)", len(c.Options))
	if len(c.Options) > 0 {
		var options []*interaction.CommandOption
		for _, o := range c.Options {
			options = append(options, o.toDiscordOption())
		}
		base.Options = options
	}
	return &subCmd{
		CommandOption: &base,
		Handler:       c.Handler,
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

// toDiscordOption turns commandOptionCreator into a interaction.CommandOption
func (o *commandOptionCreator) toDiscordOption() *interaction.CommandOption {
	var choices []*interaction.CommandOptionChoice
	for _, c := range o.Choices {
		choices = append(choices, c.toDiscordChoice())
	}
	return &interaction.CommandOption{
		Type:        o.Type,
		Name:        o.Name,
		Description: o.Description,
		Required:    o.Required,
		Choices:     choices,
	}
}

// toDiscordChoice turns commandChoiceCreator into a interaction.CommandOptionChoice (internal use of the API only)
func (c *commandChoiceCreator) toDiscordChoice() *interaction.CommandOptionChoice {
	return &interaction.CommandOptionChoice{
		Name:  c.Name,
		Value: c.Value,
	}
}
