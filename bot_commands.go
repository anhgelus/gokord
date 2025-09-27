package gokord

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/anhgelus/gokord/cmd"
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"
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
		s.Logger().Warn("update is nil, returning")
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
	for g := range s.GuildAPI().State.Guilds() {
		if g.PublicUpdatesChannelID != "" {
			changelog := fmt.Sprintf("## Nouveauté de la %s\n%s", update.Version, update.Changelog)
			_, err := s.ChannelAPI().MessageSend(g.PublicUpdatesChannelID, changelog)
			if err != nil {
				b.Logger.Error("sending changelog to guild", "error", err, "guild", g.ID)
			}
		}
	}
}

// removeCommands delete commands of InnovationCommands.Removed
func (b *Bot) removeCommands(s *discordgo.Session, update *InnovationCommands) {
	appID := s.SessionState().User().ID
	cmdRegistered, err := s.InteractionAPI().Commands(appID, "")
	if err != nil {
		b.Logger.Error("fetching slash commands", "error", err)
		return
	}
	for _, d := range update.Removed {
		id := slices.IndexFunc(cmdRegistered, func(e *interaction.Command) bool {
			return d == e.Name
		})
		if id == -1 {
			b.Logger.Warn("command not registered, so it cannot be deleted", "command", d)
			continue
		}
		c := cmdRegistered[id]
		err = s.InteractionAPI().CommandDelete(appID, "", cmdRegistered[id].ApplicationID)
		if err != nil {
			b.Logger.Error("deleting command", "error", err, "name", c.Name, "id", c.ID)
		}
	}
}

// registerCommands creates commands of InnovationCommands.Added and updates commands of InnovationCommands.Added
func (b *Bot) registerCommands(s *discordgo.Session, update *InnovationCommands) {
	var toUpdate []cmd.CommandBuilder
	guildID := ""

	// set toUpdate and guildID
	if Debug {
		gs := s.GuildAPI().State.Guilds()
		// This is ugly, but I can't do something else because Guilds returns an iter
		for g := range gs {
			guildID = g.ID
			break
		}
		if guildID == "" {
			b.Logger.Error("fetching guilds for debug", "error", fmt.Errorf("no cached guilds"))
			return
		}
		// all commands (because it is in Debug)
		toUpdate = b.Commands
	} else {
		for _, c := range append(update.Updated, update.Added...) {
			id := slices.IndexFunc(b.Commands, func(e cmd.CommandBuilder) bool {
				return c == e.GetName()
			})
			if id == -1 {
				b.Logger.Warn("impossible to find command", "command", c)
			} else {
				toUpdate = append(toUpdate, b.Commands[id])
			}
		}
	}

	// update everything needed
	appID := s.SessionState().User().ID
	o := 0
	if Debug {
		for _, cb := range toUpdate {
			registeredCommands = append(registeredCommands, cb.ApplicationCommand())
		}
		created, err := s.InteractionAPI().CommandBulkOverwrite(appID, guildID, registeredCommands)
		if err != nil {
			b.Logger.Error("registering guild commands", "error", err)
		}
		registeredCommands = created
		o = len(registeredCommands)
	} else {
		for _, cb := range toUpdate {
			c, err := s.InteractionAPI().CommandCreate(appID, "", cb.ApplicationCommand())
			if err != nil {
				b.Logger.Error("registering command", "error", err, "command", cb.GetName())
				continue
			}
			registeredCommands = append(registeredCommands, c)
			o += 1
		}
	}
	l := len(toUpdate)
	var level slog.Level
	if l != o {
		level = slog.LevelWarn
	} else {
		level = slog.LevelInfo
	}
	s.Logger().Log(context.Background(), level, "commands setups finished", "updated", o, "to update", l)
}

// setupCommandsHandlers of the Bot
func (b *Bot) setupCommandsHandlers(s *discordgo.Session) {
	if cmdMap == nil || len(cmdMap) == 0 {
		cmdMap = make(map[string]cmd.CommandHandler, len(b.Commands))
		for _, c := range b.Commands {
			b.Logger.Debug("setup handler", "command", c.GetName())
			if c.HasSub() {
				b.Logger.Debug("using general handler", "command", c.GetName())
				cmdMap[c.GetName()] = b.generalHandler
			} else {
				cmdMap[c.GetName()] = c.GetHandler()
			}
		}
		cmdMap["ping"] = pingCommand
	}
	s.EventManager().AddHandler(func(s bot.Session, i *event.InteractionCreate) {
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
		b.Logger.Error("fetching guilds for debug", "error", err)
		return
	}
	guildID := gs[0].ID
	var wg sync.WaitGroup
	for _, v := range registeredCommands {
		wg.Add(1)
		go func() {
			err = s.InteractionAPI().CommandDelete(s.SessionState().User().ID, guildID, v.ID)
			if err != nil {
				b.Logger.Error("deleting command", "error", err, "name", v.Name, "id", v.ID)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	registeredCommands = []*interaction.Command{}
}
