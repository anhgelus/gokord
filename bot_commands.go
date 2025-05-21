package gokord

import (
	"fmt"
	"github.com/anhgelus/gokord/commands"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
	"slices"
)

// updateCommands of the Bot
func (b *Bot) updateCommands(client *discordgo.Session) {
	// add ping command
	b.Commands = append(
		b.Commands,
		NewCommand("ping", "Connect the ping of the bot").SetHandler(commands.Ping),
	)
	// removing old commands and skipping already registered commands
	appID := client.State.Application.ID
	cmdRegistered, err := client.ApplicationCommands(appID, "")
	if err != nil {
		utils.SendAlert("bot.go - Fetching slash commands", err.Error())
		return
	}
	update := b.getCommandsUpdate()
	if update == nil {
		utils.SendAlert("bot.go - Checking the update", "update is nil, check the log")
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
				"bot.go - Deleting slash command", err.Error(),
				"name", c.Name,
				"id", c.ApplicationID,
			)
		}
	}
	// registering commands
	var toUpdate []*GeneralCommand
	guildID := ""

	// set toUpdate and guildID
	if Debug {
		gs, err := client.UserGuilds(1, "", "", false)
		if err != nil {
			utils.SendAlert("bot.go - Fetching guilds for debug", err.Error())
			return
		} else {
			guildID = gs[0].ID
		}
		// all commands (because it is in Debug)
		toUpdate = b.Commands
	} else {
		for _, c := range append(update.Updated, update.Added...) {
			id := slices.IndexFunc(b.Commands, func(e *GeneralCommand) bool {
				return c == e.Name
			})
			if id == -1 {
				utils.SendWarn("Impossible to find command", "name", c)
				continue
			}
			toUpdate = append(toUpdate, b.Commands[id])
		}
	}

	o := 0
	for _, c := range toUpdate {
		cmd, err := client.ApplicationCommandCreate(client.State.User.ID, guildID, c.ToCmd().ApplicationCommand)
		if err != nil {
			utils.SendAlert("bot.go - Create guild application command", err.Error(), "name", c)
			continue
		}
		registeredCommands = append(registeredCommands, cmd)
		utils.SendSuccess(fmt.Sprintf("Command %s initialized", c))
		o += 1
	}
	l := len(toUpdate)
	msg := fmt.Sprintf("%d/%d commands has been created or updated", o, l)
	if l != o {
		utils.SendWarn(msg)
	} else {
		utils.SendSuccess(msg)
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
		utils.SendAlert("bot.go - Fetching guilds for debug", err.Error())
		return
	}
	guildID := gs[0].ID
	for _, v := range registeredCommands {
		err = s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
		if err != nil {
			utils.SendAlert("bot.go - Delete application command", err.Error())
			continue
		}
	}
	registeredCommands = []*discordgo.ApplicationCommand{}
}
