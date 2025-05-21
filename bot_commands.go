package gokord

import (
	"fmt"
	"github.com/anhgelus/gokord/commands"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
	"slices"
	"sync"
)

// updateCommands of the Bot
func (b *Bot) updateCommands(s *discordgo.Session) {
	// add ping command
	b.Commands = append(
		b.Commands,
		NewCommand("ping", "Get the ping of the bot").SetHandler(commands.Ping),
	)

	update := b.getCommandsUpdate()
	if update == nil {
		utils.SendAlert("bot.go - Checking the update", "update is nil, check the log")
		return
	}

	var wg sync.WaitGroup
	// if Debug, avoid removing commands
	if !Debug {
		wg.Add(1)
		go func() {
			b.removeCommands(s, update)
			wg.Done()
		}()
	}
	wg.Add(1)
	go func() {
		b.registerCommands(s, update)
		wg.Done()
	}()
	wg.Wait()
	b.Version.UpdateBotVersion(b)
}

// removeCommands delete commands of InnovationCommands.Removed
func (b *Bot) removeCommands(s *discordgo.Session, update *InnovationCommands) {
	appID := s.State.User.ID
	cmdRegistered, err := s.ApplicationCommands(appID, "")
	if err != nil {
		utils.SendAlert("bot.go - Fetching slash commands", err.Error())
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
		err = s.ApplicationCommandDelete(appID, "", cmdRegistered[id].ApplicationID)
		if err != nil {
			utils.SendAlert(
				"bot.go - Deleting slash command", err.Error(),
				"name", c.Name,
				"id", c.ApplicationID,
			)
		}
	}
}

// registerCommands creates commands of InnovationCommands.Added and updates commands of InnovationCommands.Added
func (b *Bot) registerCommands(s *discordgo.Session, update *InnovationCommands) {
	var toUpdate []CommandBuilder
	guildID := ""

	// set toUpdate and guildID
	if Debug {
		gs, err := s.UserGuilds(1, "", "", false)
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
			id := slices.IndexFunc(b.Commands, func(e CommandBuilder) bool {
				return c == e.toCreator().Name
			})
			if id == -1 {
				utils.SendWarn("Impossible to find command", "name", c)
				continue
			}
			toUpdate = append(toUpdate, b.Commands[id])
		}
	}

	// update everything needed
	appID := s.State.User.ID
	o := 0
	for _, cb := range toUpdate {
		c, err := s.ApplicationCommandCreate(appID, guildID, cb.toCreator().ToCmd().ApplicationCommand)
		if err != nil {
			utils.SendAlert("bot.go - Create guild application command", err.Error(), "name", cb)
			continue
		}
		registeredCommands = append(registeredCommands, c)
		utils.SendSuccess(fmt.Sprintf("Command %s initialized", cb.toCreator().Name))
		o += 1
	}
	l := len(toUpdate)
	msg := fmt.Sprintf("%d/%d commands has been created or updated", o, l)
	if l != o {
		utils.SendWarn(msg)
	} else {
		utils.SendSuccess(msg)
	}
}

// setupCommandsHandlers of the Bot
func (b *Bot) setupCommandsHandlers(s *discordgo.Session) {
	if len(cmdMap) == 0 {
		for _, cb := range b.Commands {
			c := cb.toCreator()
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
