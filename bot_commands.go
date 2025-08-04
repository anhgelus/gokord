package gokord

import (
	"fmt"
	"github.com/anhgelus/gokord/cmd"
	"github.com/anhgelus/gokord/logger"
	"github.com/bwmarrin/discordgo"
	"slices"
	"sync"
)

var cmdMap map[string]cmd.CommandHandler = nil

// updateCommands of the Bot
func (b *Bot) updateCommands(s *discordgo.Session) {
	// add ping command
	b.Commands = append(
		b.Commands,
		cmd.New("ping", "Get the ping of the bot").
			SetHandler(pingCommand).
			AddContext(discordgo.InteractionContextGuild).
			AddContext(discordgo.InteractionContextBotDM).
			AddContext(discordgo.InteractionContextPrivateChannel).
			AddIntegrationType(discordgo.ApplicationIntegrationGuildInstall).
			AddIntegrationType(discordgo.ApplicationIntegrationUserInstall),
	)

	update := b.getCommandsUpdate()
	if update == nil {
		logger.Alert("bot.go - Checking the update", "update is nil, check the log")
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
		logger.Alert("bot.go - Fetching slash commands", err.Error())
		return
	}
	for _, d := range update.Removed {
		id := slices.IndexFunc(cmdRegistered, func(e *discordgo.ApplicationCommand) bool {
			return d == e.Name
		})
		if id == -1 {
			logger.Warn("Command not registered cannot be deleted", "name", d)
			continue
		}
		c := cmdRegistered[id]
		err = s.ApplicationCommandDelete(appID, "", cmdRegistered[id].ApplicationID)
		if err != nil {
			logger.Alert(
				"bot.go - Deleting slash command", err.Error(),
				"name", c.Name,
				"id", c.ApplicationID,
			)
		}
	}
}

// registerCommands creates commands of InnovationCommands.Added and updates commands of InnovationCommands.Added
func (b *Bot) registerCommands(s *discordgo.Session, update *InnovationCommands) {
	var toUpdate []cmd.CommandBuilder
	guildID := ""

	// set toUpdate and guildID
	if Debug {
		gs, err := s.UserGuilds(1, "", "", false)
		if err != nil {
			logger.Alert("bot.go - Fetching guilds for debug", err.Error())
			return
		} else {
			guildID = gs[0].ID
		}
		// all commands (because it is in Debug)
		toUpdate = b.Commands
	} else {
		for _, c := range append(update.Updated, update.Added...) {
			id := slices.IndexFunc(b.Commands, func(e cmd.CommandBuilder) bool {
				return c == e.GetName()
			})
			if id == -1 {
				logger.Warn("Impossible to find command", "name", c)
				continue
			}
			toUpdate = append(toUpdate, b.Commands[id])
		}
	}

	// update everything needed
	appID := s.State.User.ID
	o := 0
	for _, cb := range toUpdate {
		c, err := s.ApplicationCommandCreate(appID, guildID, cb.ApplicationCommand())
		if err != nil {
			logger.Alert("bot.go - Create guild application command", err.Error(), "name", cb.GetName())
			continue
		}
		registeredCommands = append(registeredCommands, c)
		logger.Success(fmt.Sprintf("Command %s initialized", cb.GetName()))
		o += 1
	}
	l := len(toUpdate)
	msg := fmt.Sprintf("%d/%d commands has been created or updated", o, l)
	if l != o {
		logger.Warn(msg)
	} else {
		logger.Success(msg)
	}
}

// setupCommandsHandlers of the Bot
func (b *Bot) setupCommandsHandlers(s *discordgo.Session) {
	if cmdMap == nil || len(cmdMap) == 0 {
		cmdMap = make(map[string]cmd.CommandHandler, len(b.Commands))
		for _, c := range b.Commands {
			logger.Debug("Setup handler", "command", c.GetName())
			if c.HasSub() {
				logger.Debug("Using general handler", "command", c.GetName())
				cmdMap[c.GetName()] = b.generalHandler
			} else {
				cmdMap[c.GetName()] = c.GetHandler()
			}
		}
		cmdMap["ping"] = pingCommand
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		if h, ok := cmdMap[i.ApplicationCommandData().Name]; ok {
			resp := cmd.NewResponseBuilder(s, i)
			optMap := cmd.GenerateOptionMap(i)
			h(s, i, optMap, resp)
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
		logger.Alert("bot.go - Fetching guilds for debug", err.Error())
		return
	}
	guildID := gs[0].ID
	var wg sync.WaitGroup
	for _, v := range registeredCommands {
		wg.Add(1)
		go func() {
			err = s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
			if err != nil {
				logger.Alert("bot.go - Delete application command", err.Error())
			}
			wg.Done()
		}()
	}
	wg.Wait()
	registeredCommands = []*discordgo.ApplicationCommand{}
}
