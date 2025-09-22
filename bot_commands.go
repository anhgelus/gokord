package gokord

import (
	"fmt"
	"slices"
	"sync"

	"github.com/anhgelus/gokord/cmd"
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/interaction"
	"github.com/nyttikord/gokord/logger"
)

var cmdMap map[string]cmd.CommandHandler = nil

// updateCommands of the Bot
func (b *Bot) updateCommands(s *discordgo.Session) {
	// add ping command
	b.Commands = append(
		b.Commands,
		cmd.New("ping", "Get the ping of the bot").
			SetHandler(pingCommand).
			AddContext(types.InteractionContextGuild).
			AddContext(types.InteractionContextBotDM).
			AddContext(types.InteractionContextPrivateChannel).
			AddIntegrationType(types.IntegrationInstallGuild).
			AddIntegrationType(types.IntegrationInstallUser),
	)

	update, do := b.getCommandsUpdate()
	if !do {
		return
	}
	if update == nil {
		s.LogWarn("Update is nil, returning.")
		return
	}

	var wg sync.WaitGroup
	// if Debug, avoid removing commands
	if !Debug {
		wg.Add(1)
		go func() {
			b.removeCommands(s, update.Commands)
			wg.Done()
		}()
	}
	wg.Add(1)
	go func() {
		b.registerCommands(s, update.Commands)
		wg.Done()
	}()
	wg.Wait()
	b.Version.UpdateBotVersion(b)
	// sending changelog to guilds
	if update.Changelog == "" {
		return
	}
	for _, g := range s.State.Guilds {
		if g.PublicUpdatesChannelID != "" {
			changelog := fmt.Sprintf("## Nouveauté de la %s\n%s", update.Version, update.Changelog)
			_, err := s.ChannelAPI().MessageSend(g.PublicUpdatesChannelID, changelog)
			if err != nil {
				b.LogError(err, "sending changelog to guild %s", g.ID)
			}
		}
	}
}

// removeCommands delete commands of InnovationCommands.Removed
func (b *Bot) removeCommands(s *discordgo.Session, update *InnovationCommands) {
	appID := s.State.User.ID
	cmdRegistered, err := s.InteractionAPI().Commands(appID, "")
	if err != nil {
		s.LogError(err, "fetching slash commands")
		return
	}
	for _, d := range update.Removed {
		id := slices.IndexFunc(cmdRegistered, func(e *interaction.Command) bool {
			return d == e.Name
		})
		if id == -1 {
			s.LogWarn("Command %s not registered, so it cannot be deleted", d)
			continue
		}
		c := cmdRegistered[id]
		err = s.InteractionAPI().CommandDelete(appID, "", cmdRegistered[id].ApplicationID)
		if err != nil {
			s.LogError(err, "deleting command %s (%s)", c.Name, c.ID)
		}
	}
}

// registerCommands creates commands of InnovationCommands.Added and updates commands of InnovationCommands.Added
func (b *Bot) registerCommands(s *discordgo.Session, update *InnovationCommands) {
	var toUpdate []cmd.CommandBuilder
	guildID := ""

	// set toUpdate and guildID
	if Debug {
		gs, err := s.GuildAPI().UserGuilds(1, "", "", false)
		if err != nil {
			s.LogError(err, "fetching guilds for debug")
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
				s.LogWarn("Impossible to find command %s", c)
				continue
			}
			toUpdate = append(toUpdate, b.Commands[id])
		}
	}

	// update everything needed
	appID := s.State.User.ID
	o := 0
	for _, cb := range toUpdate {
		c, err := s.InteractionAPI().CommandCreate(appID, guildID, cb.ApplicationCommand())
		if err != nil {
			s.LogError(err, "creating guild command %s", cb.GetName())
			continue
		}
		registeredCommands = append(registeredCommands, c)
		s.LogInfo("Command %s initialized", cb.GetName())
		o += 1
	}
	l := len(toUpdate)
	var level logger.Level
	if l != o {
		level = logger.LevelWarn
	} else {
		level = logger.LevelInfo
	}
	s.Log(level, 0, "%d/%d commands has been created or updated", o, l)
}

// setupCommandsHandlers of the Bot
func (b *Bot) setupCommandsHandlers(s *discordgo.Session) {
	if cmdMap == nil || len(cmdMap) == 0 {
		cmdMap = make(map[string]cmd.CommandHandler, len(b.Commands))
		for _, c := range b.Commands {
			s.LogDebug("Setup handler for command %s", c.GetName())
			if c.HasSub() {
				b.LogDebug("Using general handler for command %s", c.GetName())
				cmdMap[c.GetName()] = b.generalHandler
			} else {
				cmdMap[c.GetName()] = c.GetHandler()
			}
		}
		cmdMap["ping"] = pingCommand
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != types.InteractionApplicationCommand {
			return
		}
		if h, ok := cmdMap[i.CommandData().Name]; ok {
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
	gs, err := s.GuildAPI().UserGuilds(1, "", "", false)
	if err != nil {
		s.LogError(err, "fetching guilds for debug")
		return
	}
	guildID := gs[0].ID
	var wg sync.WaitGroup
	for _, v := range registeredCommands {
		wg.Add(1)
		go func() {
			err = s.InteractionAPI().CommandDelete(s.State.User.ID, guildID, v.ID)
			if err != nil {
				s.LogError(err, "deleting command %s (%s)", v.Name, v.ID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	registeredCommands = []*interaction.Command{}
}
