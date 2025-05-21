package gokord

import (
	"errors"
	"fmt"
	"github.com/anhgelus/gokord/commands"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
	"math/rand/v2"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"
)

var (
	Debug = true

	ErrBadStatusType     = errors.New("bad status type, please use the constant")
	ErrStatusUrlNotFound = errors.New("status url not found")

	registeredCommands []*discordgo.ApplicationCommand
)

const (
	GameStatus      = 0
	WatchStatus     = 1
	StreamingStatus = 2
	ListeningStatus = 3

	AdminPermission int64 = discordgo.PermissionManageServer // AdminPermission of the command
)

// Bot is the representation of a discord bot
type Bot struct {
	Token    string            // Token of the Bot
	Status   []*Status         // Status of the Bot
	Commands []*GeneralCommand // Commands of the Bot, use NewCommand to create easily a new command
	//Handlers []interface{} // Handlers of the Bot
	AfterInit   func(s *discordgo.Session) // AfterInit is called after the initialization process of the Bot
	Version     *Version
	Innovations []*Innovation
	Name        string
}

// Status contains all required information for updating the status
type Status struct {
	Type    int    // Type of the Status (use GameStatus or WatchStatus or StreamingStatus or ListeningStatus)
	Content string // Content of the Status
	Url     string // Url of the StreamingStatus
}

// Start the Bot (blocking instruction)
func (b *Bot) Start() {
	dg, err := discordgo.New("Bot " + b.Token) // New connection to the discord API with bot token
	if err != nil {
		utils.SendAlert("bot.go - Token", err.Error())
		return
	}

	err = dg.Open() // Starts the bot
	if err != nil {
		utils.SendAlert("bot.go - Start", err.Error())
	}
	dg.Identify.Intents = discordgo.IntentsAll

	var wg sync.WaitGroup
	st1 := time.Now()
	// register commands
	wg.Add(1)
	go func() {
		b.updateCommands(dg)
		utils.SendSuccess("Commands updated")
		wg.Done()
	}()
	b.setupCommandsHandlers(dg)

	//for h := range b.Handlers {
	//	dg.AddHandler(h)
	//}

	// do after init (mainly used to register handlers)
	b.AfterInit(dg)

	// wait until all setup goroutines are finished
	wg.Wait()
	st2 := time.Now()
	delta := float64(st2.Unix() - st1.Unix())
	to := dg.Client.Timeout.Seconds()
	// if the setup was faster than the http client timeout, wait
	if delta < dg.Client.Timeout.Seconds() {
		time.Sleep(time.Duration(to-delta) * time.Second)
	}
	utils.SendSuccess(fmt.Sprintf("Bot started as %s", dg.State.User.Username))
	utils.NewTimer(30*time.Second, func(stop chan<- interface{}) {
		if b.Status == nil {
			stop <- struct{}{}
			return
		}
		l := len(b.Status)
		r := rand.New(rand.NewPCG(uint64(time.Now().Unix()), uint64(l))).UintN(uint(l))
		s := b.Status[r]
		if s.Type == GameStatus {
			err = dg.UpdateGameStatus(0, s.Content)
		} else if s.Type == WatchStatus {
			err = dg.UpdateWatchStatus(0, s.Content)
		} else if s.Type == StreamingStatus {
			if s.Url == "" {
				err = ErrStatusUrlNotFound
			} else {
				err = dg.UpdateStreamingStatus(0, s.Content, s.Url)
			}
		} else if s.Type == ListeningStatus {
			err = dg.UpdateListeningStatus(s.Content)
		} else {
			err = ErrBadStatusType
		}
		if err != nil {
			utils.SendAlert("bot.go - Update status", err.Error())
			err = nil
		}
	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	utils.SendSuccess("Bot shutting down")

	if Debug {
		utils.SendDebug("Unregistering local commands")
		b.unregisterGuildCommands(dg)
	}

	err = dg.Close() // Bot Shutdown
	if err != nil {
		utils.SendAlert("bot.go - Shutdown", err.Error())
	}

	utils.SendSuccess("Bot shut down")
}

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
	update := GetCommandsUpdate(b)
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
				"bot.go - Deleting slash command",
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
			utils.SendAlert("bot.go - Fetching guilds for debug", err.Error())
			return
		} else {
			guildID = gs[0].ID
		}
	}
	var registeredCommands []*discordgo.ApplicationCommand
	o := 0
	for _, c := range append(update.Updated, update.Added...) {
		o += 1
		id := slices.IndexFunc(b.Commands, func(e *GeneralCommand) bool {
			return c == e.Name
		})
		v := b.Commands[id]
		cmd, err := client.ApplicationCommandCreate(client.State.User.ID, guildID, v.ToCmd().ApplicationCommand)
		if err != nil {
			utils.SendAlert("bot.go - Create application command", err.Error(), "name", c)
			continue
		}
		registeredCommands = append(registeredCommands, cmd)
		utils.SendSuccess(fmt.Sprintf("Command %s initialized", c))
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
