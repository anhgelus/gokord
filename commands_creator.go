package gokord

import (
	"github.com/bwmarrin/discordgo"
)

// Cmd is a discordgo.ApplicationCommand + its handler
//
// Use AdminPermission to set the admin permission
type Cmd struct {
	*discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate) // Handler called
	Subs    []*SimpleSubCmd
}

// SubCmd is for the internal use of the API
type SubCmd struct {
	*discordgo.ApplicationCommandOption
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate) // Handler called
}

// SimpleSubCmd is for the internal use of the API
type SimpleSubCmd struct {
	Name    string
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate) // Handler called
}

// CommandCreator represents a generic command
type CommandCreator struct {
	HasSub      bool
	IsSub       bool
	Name        string
	Permission  *int64
	CanDM       bool
	Description string
	Options     []*CommandOptionCreator
	Subs        []*CommandCreator
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate) // Handler called
}

// CommandOptionCreator represents a generic option of CommandCreator
type CommandOptionCreator struct {
	Type        discordgo.ApplicationCommandOptionType
	Name        string
	Description string
	Required    bool
	Choices     []*CommandChoiceCreator
}

// CommandChoiceCreator represents a generic choice of CommandOptionCreator
type CommandChoiceCreator struct {
	Name  string
	Value interface{}
}

var (
	cmdMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

// ToSimple turns SubCmd into a SimpleSubCmd
func (s *SubCmd) ToSimple() *SimpleSubCmd {
	return &SimpleSubCmd{
		Name:    s.Name,
		Handler: s.Handler,
	}
}

// NewCommand creates a CommandCreator
func NewCommand(name string, description string) *CommandCreator {
	return &CommandCreator{
		HasSub:      false,
		IsSub:       false,
		Name:        name,
		CanDM:       false,
		Description: description,
		Options:     []*CommandOptionCreator{},
		Subs:        []*CommandCreator{},
	}
}

// SetHandler of the CommandCreator (if CommandCreator contains subcommand, it will never be called)
func (c *CommandCreator) SetHandler(handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) *CommandCreator {
	c.Handler = handler
	return c
}

// ContainsSub makes the CommandCreator able to contain subcommands
func (c *CommandCreator) ContainsSub() *CommandCreator {
	c.HasSub = true
	c.Options = nil
	return c
}

// AddSub to the CommandCreator (also call ContainsSub)
func (c *CommandCreator) AddSub(s *CommandCreator) *CommandCreator {
	c.ContainsSub()
	s.IsSub = true
	c.Subs = append(c.Subs, s)
	return c
}

// HasOption makes the CommandCreator able to contain CommandOptionCreator
func (c *CommandCreator) HasOption() *CommandCreator {
	c.HasSub = false
	c.Subs = nil
	return c
}

// AddOption to the CommandCreator (also call HasOption)
func (c *CommandCreator) AddOption(s *CommandOptionCreator) *CommandCreator {
	c.HasOption()
	c.Options = append(c.Options, s)
	return c
}

// DM makes the CommandCreator used in DM
func (c *CommandCreator) DM() *CommandCreator {
	c.CanDM = true
	return c
}

// SetPermission of the CommandCreator
func (c *CommandCreator) SetPermission(p *int64) *CommandCreator {
	c.Permission = p
	return c
}

// Is returns true if the CommandCreator is approximately the same as *discordgo.ApplicationCommand
func (c *CommandCreator) Is(cmd *discordgo.ApplicationCommand) bool {
	return cmd.DefaultMemberPermissions == c.Permission &&
		cmd.Name == c.Name &&
		cmd.Description == c.Description &&
		len(cmd.Options) == len(c.Options)
}

// ToCmd turns CommandCreator into a Cmd (internal use of the API only)
func (c *CommandCreator) ToCmd() *Cmd {
	base := discordgo.ApplicationCommand{
		Type:         discordgo.ChatApplicationCommand,
		Name:         c.Name,
		DMPermission: &c.CanDM,
		Description:  c.Description,
	}
	if c.Permission != nil {
		base.DefaultMemberPermissions = c.Permission
	}
	if !c.HasSub {
		var options []*discordgo.ApplicationCommandOption
		for _, o := range c.Options {
			options = append(options, o.ToDiscordOption())
		}
		base.Options = options
		return &Cmd{
			ApplicationCommand: &base,
			Handler:            c.Handler,
		}
	}
	var subsCmd []*SimpleSubCmd
	var subs []*discordgo.ApplicationCommandOption
	for _, s := range c.Subs {
		sub := s.ToSubCmd()
		subsCmd = append(subsCmd, sub.ToSimple())
		subs = append(subs, sub.ApplicationCommandOption)
	}
	base.Options = subs
	return &Cmd{
		ApplicationCommand: &base,
		Handler:            c.Handler,
		Subs:               subsCmd,
	}
}

// ToSubCmd turns CommandCreator into a SubCmd (internal use of the API only)
func (c *CommandCreator) ToSubCmd() *SubCmd {
	base := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        c.Name,
		Description: c.Description,
	}
	if c.Options != nil {
		var options []*discordgo.ApplicationCommandOption
		for _, o := range c.Options {
			options = append(options, o.ToDiscordOption())
		}
		base.Options = options
	}
	return &SubCmd{
		ApplicationCommandOption: &base,
		Handler:                  c.Handler,
	}
}

// NewOption creates a new CommandOptionCreator
func NewOption(t discordgo.ApplicationCommandOptionType, name string, description string) *CommandOptionCreator {
	return &CommandOptionCreator{
		Type:        t,
		Name:        name,
		Description: description,
		Required:    false,
		Choices:     []*CommandChoiceCreator{},
	}
}

// IsRequired informs that the CommandOptionCreator is required
func (o *CommandOptionCreator) IsRequired() *CommandOptionCreator {
	o.Required = true
	return o
}

// AddChoice to the CommandOptionCreator
func (o *CommandOptionCreator) AddChoice(c *CommandChoiceCreator) *CommandOptionCreator {
	o.Required = true
	o.Choices = append(o.Choices, c)
	return o
}

// ToDiscordOption turns CommandOptionCreator into a discordgo.ApplicationCommandOption (internal use of the API only)
func (o *CommandOptionCreator) ToDiscordOption() *discordgo.ApplicationCommandOption {
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

// NewChoice creates a new choice for CommandOptionCreator
func NewChoice(name string, value interface{}) *CommandChoiceCreator {
	return &CommandChoiceCreator{
		Name:  name,
		Value: value,
	}
}

// ToDiscordChoice turns CommandChoiceCreator into a discordgo.ApplicationCommandOptionChoice (internal use of the API only)
func (c *CommandChoiceCreator) ToDiscordChoice() *discordgo.ApplicationCommandOptionChoice {
	return &discordgo.ApplicationCommandOptionChoice{
		Name:  c.Name,
		Value: c.Value,
	}
}
