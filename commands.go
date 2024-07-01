package gokord

import (
	"fmt"
	"github.com/anhgelus/gokord/commands"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
	"slices"
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

// GeneralCommand represents a generic command
type GeneralCommand struct {
	HasSub      bool
	IsSub       bool
	Name        string
	Permission  *int64
	CanDM       bool
	Description string
	Options     []*GCommandOption
	Subs        []*GeneralCommand
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate) // Handler called
}

// GCommandOption represents a generic option of GeneralCommand
type GCommandOption struct {
	Type        discordgo.ApplicationCommandOptionType
	Name        string
	Description string
	Required    bool
	Choices     []*GCommandChoice
}

// GCommandChoice represents a generic choice of GCommandOption
type GCommandChoice struct {
	Name  string
	Value interface{}
}

var (
	cmdMap             = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	registeredCommands []*discordgo.ApplicationCommand
)

// ToSimple turns SubCmd into a SimpleSubCmd
func (s *SubCmd) ToSimple() *SimpleSubCmd {
	return &SimpleSubCmd{
		Name:    s.Name,
		Handler: s.Handler,
	}
}

// NewCommand creates a GeneralCommand
func NewCommand(name string, description string) *GeneralCommand {
	return &GeneralCommand{
		HasSub:      false,
		IsSub:       false,
		Name:        name,
		CanDM:       false,
		Description: description,
		Options:     []*GCommandOption{},
		Subs:        []*GeneralCommand{},
	}
}

// SetHandler of the GeneralCommand (if GeneralCommand contains subcommand, it will never be called)
func (c *GeneralCommand) SetHandler(handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) *GeneralCommand {
	c.Handler = handler
	return c
}

// ContainsSub makes the GeneralCommand able to contain subcommands
func (c *GeneralCommand) ContainsSub() *GeneralCommand {
	c.HasSub = true
	c.Options = nil
	return c
}

// AddSub to the GeneralCommand (also call ContainsSub)
func (c *GeneralCommand) AddSub(s *GeneralCommand) *GeneralCommand {
	c.ContainsSub()
	s.IsSub = true
	c.Subs = append(c.Subs, s)
	return c
}

// HasOption makes the GeneralCommand able to contains GCommandOption
func (c *GeneralCommand) HasOption() *GeneralCommand {
	c.HasSub = false
	c.Subs = nil
	return c
}

// AddOption to the GeneralCommand (also call HasOption)
func (c *GeneralCommand) AddOption(s *GCommandOption) *GeneralCommand {
	c.HasOption()
	c.Options = append(c.Options, s)
	return c
}

// DM makes the GeneralCommand used in DM
func (c *GeneralCommand) DM() *GeneralCommand {
	c.CanDM = true
	return c
}

// SetPermission of the GeneralCommand
func (c *GeneralCommand) SetPermission(p *int64) *GeneralCommand {
	c.Permission = p
	return c
}

// Is returns true if the GeneralCommand is approximately the same as *discordgo.ApplicationCommand
func (c *GeneralCommand) Is(cmd *discordgo.ApplicationCommand) bool {
	return cmd.DefaultMemberPermissions == c.Permission &&
		cmd.Name == c.Name &&
		cmd.Description == c.Description &&
		len(cmd.Options) == len(c.Options)
}

// ToCmd turns GeneralCommand into a Cmd (internal use of the API only)
func (c *GeneralCommand) ToCmd() *Cmd {
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

// ToSubCmd turns GeneralCommand into a SubCmd (internal use of the API only)
func (c *GeneralCommand) ToSubCmd() *SubCmd {
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

// NewOption creates a new GCommandOption
func NewOption(t discordgo.ApplicationCommandOptionType, name string, description string) *GCommandOption {
	return &GCommandOption{
		Type:        t,
		Name:        name,
		Description: description,
		Required:    false,
		Choices:     []*GCommandChoice{},
	}
}

// IsRequired informs that the GCommandOption is required
func (o *GCommandOption) IsRequired() *GCommandOption {
	o.Required = true
	return o
}

// AddChoice to the GCommandOption
func (o *GCommandOption) AddChoice(c *GCommandChoice) *GCommandOption {
	o.Required = true
	o.Choices = append(o.Choices, c)
	return o
}

// ToDiscordOption turns GCommandOption into a discordgo.ApplicationCommandOption (internal use of the API only)
func (o *GCommandOption) ToDiscordOption() *discordgo.ApplicationCommandOption {
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

// NewChoice creates a new choice for GCommandOption
func NewChoice(name string, value interface{}) *GCommandChoice {
	return &GCommandChoice{
		Name:  name,
		Value: value,
	}
}

// ToDiscordChoice turns GCommandChoice into a discordgo.ApplicationCommandOptionChoice (internal use of the API only)
func (c *GCommandChoice) ToDiscordChoice() *discordgo.ApplicationCommandOptionChoice {
	return &discordgo.ApplicationCommandOptionChoice{
		Name:  c.Name,
		Value: c.Value,
	}
}

// updateCommands of the Bot
func (b *Bot) updateCommands(client *discordgo.Session) {
	// add ping command
	b.Commands = append(
		b.Commands,
		NewCommand("ping", "Get the ping of the bot").SetHandler(commands.Ping),
	)
	// removing old commands and skipping already registered commands
	appID := client.State.Application.ID
	cmdRegistered, err := client.ApplicationCommands(appID, "")
	if err != nil {
		utils.SendAlert("commands.go - Fetching slash commands", err.Error())
		return
	}
	update := GetCommandsUpdate(b)
	if update == nil {
		utils.SendAlert("commands.go - Checking the update", "update is nil, check the log")
		return
	}
	for _, d := range update.Removed {
		id := slices.IndexFunc(cmdRegistered, func(e *discordgo.ApplicationCommand) bool {
			return d == e.Name
		})
		if id == -1 {
			utils.SendWarn("Command not registered cannot be deleted", "name", d)
			continue
		}
		c := cmdRegistered[id]
		err = client.ApplicationCommandDelete(appID, "", cmdRegistered[id].ApplicationID)
		if err != nil {
			utils.SendAlert(
				"commands.go - Deleting slash command",
				err.Error(),
				"name",
				c.Name,
				"id",
				c.ApplicationID,
			)
		}
	}
	// registering commands
	guildID := ""
	if Debug {
		gs, err := client.UserGuilds(1, "", "", false)
		if err != nil {
			utils.SendAlert("commands.go - Fetching guilds for debug", err.Error())
			return
		} else {
			guildID = gs[0].ID
		}
	}
	var registeredCommands []*discordgo.ApplicationCommand
	o := 0
	for _, c := range append(update.Updated, update.Added...) {
		id := slices.IndexFunc(b.Commands, func(e *GeneralCommand) bool {
			return c == e.Name
		})
		v := b.Commands[id]
		cmd, err := client.ApplicationCommandCreate(client.State.User.ID, guildID, v.ToCmd().ApplicationCommand)
		if err != nil {
			utils.SendAlert("commands.go - Create application command", err.Error(), "name", c)
			continue
		}
		registeredCommands = append(registeredCommands, cmd)
		utils.SendSuccess(fmt.Sprintf("Command %s initialized", c))
		o += 1
	}
	l := len(registeredCommands)
	if l != o {
		utils.SendWarn(fmt.Sprintf("%d/%d commands has been created or updated", o, l))
	} else {
		utils.SendSuccess(fmt.Sprintf("%d/%d commands has been created or updated", o, l))
	}
	b.Version.UpdateBotVersion(b)
}

// setupCommandsHandlers of the Bot
func (b *Bot) setupCommandsHandlers(s *discordgo.Session) {
	if len(cmdMap) == 0 {
		for _, c := range b.Commands {
			utils.SendDebug("Setup handler", "command", c.Name)
			if c.Subs != nil {
				utils.SendDebug("Using general handler", "command", c.Name)
				cmdMap[c.Name] = b.generalHandler
			} else {
				cmdMap[c.Name] = c.Handler
			}
		}
		cmdMap["ping"] = commands.Ping
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := cmdMap[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

// unregisterGuildCommands used to unregister commands after closing the bot (Debug = true only)
func (b *Bot) unregisterGuildCommands(s *discordgo.Session) {
	if !Debug {
		return
	}
	gs, err := s.UserGuilds(1, "", "", false)
	if err != nil {
		utils.SendAlert("commands.go - Fetching guilds for debug", err.Error())
		return
	}
	guildID := gs[0].ID
	for _, v := range registeredCommands {
		err = s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
		if err != nil {
			utils.SendAlert("commands.go - Delete application command", err.Error())
			continue
		}
	}
	registeredCommands = []*discordgo.ApplicationCommand{}
}
