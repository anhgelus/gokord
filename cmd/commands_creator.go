package cmd

import (
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
)

// cmd is a discordgo.ApplicationCommand + its handler
//
// Use AdminPermission to set the admin permission
type cmd struct {
	*discordgo.ApplicationCommand
	Handler CommandHandler // Handler called
	Subs    []*simpleSubCmd
}

// subCmd is for the internal use of the API
type subCmd struct {
	*discordgo.ApplicationCommandOption
	Handler CommandHandler // Handler called
}

// simpleSubCmd is for the internal use of the API
type simpleSubCmd struct {
	Name    string
	Handler CommandHandler // Handler called
}

// commandCreator represents a generic command
type commandCreator struct {
	HasSub           bool
	IsSub            bool
	Name             string
	Permission       *int64
	Contexts         []discordgo.InteractionContextType
	IntegrationTypes []discordgo.ApplicationIntegrationType
	Description      string
	Options          []*commandOptionCreator
	Subs             []*commandCreator
	Handler          CommandHandler // Handler called
}

// commandOptionCreator represents a generic option of commandCreator
type commandOptionCreator struct {
	Type        discordgo.ApplicationCommandOptionType
	Name        string
	Description string
	Required    bool
	Choices     []*commandChoiceCreator
}

// commandChoiceCreator represents a generic choice of commandOptionCreator
type commandChoiceCreator struct {
	Name  string
	Value interface{}
}

// ToSimple turns subCmd into a simpleSubCmd
func (s *subCmd) ToSimple() *simpleSubCmd {
	return &simpleSubCmd{
		Name:    s.Name,
		Handler: s.Handler,
	}
}

// SetHandler of the commandCreator (if commandCreator contains subcommand, it will never be called)
func (c *commandCreator) SetHandler(handler CommandHandler) *commandCreator {
	c.Handler = handler
	return c
}

// ContainsSub makes the commandCreator able to contain subcommands
func (c *commandCreator) ContainsSub() *commandCreator {
	c.HasSub = true
	c.Options = []*commandOptionCreator{}
	return c
}

// AddSub to the commandCreator (also call ContainsSub)
func (c *commandCreator) AddSub(s *commandCreator) *commandCreator {
	c.ContainsSub()
	s.IsSub = true
	c.Subs = append(c.Subs, s)
	return c
}

// HasOption makes the commandCreator able to contain commandOptionCreator
func (c *commandCreator) HasOption() *commandCreator {
	c.HasSub = false
	c.Subs = []*commandCreator{}
	return c
}

// AddOption to the commandCreator (also call HasOption)
func (c *commandCreator) AddOption(s *commandOptionCreator) *commandCreator {
	c.HasOption()
	c.Options = append(c.Options, s)
	return c
}

// AddContext to the command.
// If commandCreator.Contexts is empty, discordgo.InteractionContextGuild will be added automatically
func (c *commandCreator) AddContext(ctx discordgo.InteractionContextType) *commandCreator {
	if c.Contexts == nil {
		c.Contexts = []discordgo.InteractionContextType{}
	}
	c.Contexts = append(c.Contexts, ctx)
	return c
}

func (c *commandCreator) AddIntegrationType(it discordgo.ApplicationIntegrationType) *commandCreator {
	if c.IntegrationTypes == nil {
		c.IntegrationTypes = []discordgo.ApplicationIntegrationType{}
	}
	c.IntegrationTypes = append(c.IntegrationTypes, it)
	return c
}

// SetPermission of the commandCreator
func (c *commandCreator) SetPermission(p *int64) *commandCreator {
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

// ToCmd turns commandCreator into a cmd
func (c *commandCreator) ToCmd() *cmd {
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
	if !c.HasSub {
		var options []*discordgo.ApplicationCommandOption
		for _, o := range c.Options {
			options = append(options, o.ToDiscordOption())
		}
		base.Options = options
		return &cmd{
			ApplicationCommand: &base,
			Handler:            c.Handler,
		}
	}
	var subsCmd []*simpleSubCmd
	var subs []*discordgo.ApplicationCommandOption
	for _, s := range c.Subs {
		sub := s.ToSubCmd()
		subsCmd = append(subsCmd, sub.ToSimple())
		subs = append(subs, sub.ApplicationCommandOption)
	}
	base.Options = subs
	return &cmd{
		ApplicationCommand: &base,
		Handler:            c.Handler,
		Subs:               subsCmd,
	}
}

// ToSubCmd turns commandCreator into a subCmd
func (c *commandCreator) ToSubCmd() *subCmd {
	base := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        c.Name,
		Description: c.Description,
	}
	utils.SendDebug("Subcommand creation", "name", c.Name, "len(options)", len(c.Options))
	if len(c.Options) > 0 {
		var options []*discordgo.ApplicationCommandOption
		for _, o := range c.Options {
			options = append(options, o.ToDiscordOption())
		}
		base.Options = options
	}
	return &subCmd{
		ApplicationCommandOption: &base,
		Handler:                  c.Handler,
	}
}

// IsRequired informs that the commandOptionCreator is required
func (o *commandOptionCreator) IsRequired() *commandOptionCreator {
	o.Required = true
	return o
}

// AddChoice to the commandOptionCreator
func (o *commandOptionCreator) AddChoice(c *commandChoiceCreator) *commandOptionCreator {
	o.Required = true
	o.Choices = append(o.Choices, c)
	return o
}

// ToDiscordOption turns commandOptionCreator into a discordgo.ApplicationCommandOption
func (o *commandOptionCreator) ToDiscordOption() *discordgo.ApplicationCommandOption {
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, c := range o.Choices {
		choices = append(choices, c.ToDiscordChoice())
	}
	return &discordgo.ApplicationCommandOption{
		Type:        o.Type,
		Name:        o.Name,
		Description: o.Description,
		Required:    o.Required,
		Choices:     choices,
	}
}

// ToDiscordChoice turns commandChoiceCreator into a discordgo.ApplicationCommandOptionChoice (internal use of the API only)
func (c *commandChoiceCreator) ToDiscordChoice() *discordgo.ApplicationCommandOptionChoice {
	return &discordgo.ApplicationCommandOptionChoice{
		Name:  c.Name,
		Value: c.Value,
	}
}
