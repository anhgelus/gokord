package gokord

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/anhgelus/gokord/cmd"
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"
	"github.com/nyttikord/gokord/logger"
)

var (
	Debug = true

	ErrBadStatusType     = errors.New("bad status type, please use the constant")
	ErrStatusUrlNotFound = errors.New("status url not found")

	registeredCommands []*interaction.Command
)

type StatusType int

const (
	GameStatus      StatusType = 0
	WatchStatus     StatusType = 1
	StreamingStatus StatusType = 2
	ListeningStatus StatusType = 3

	AdminPermission int64 = discord.PermissionManageGuild // AdminPermission of the command
)

// Bot is the representation of a discord bot
type Bot struct {
	logger.Logger
	Token       string                     // Token of the Bot
	Status      []*Status                  // Status of the Bot
	Commands    []cmd.CommandBuilder       // Commands of the Bot, use New to create easily a new command
	handlers    []any                      // handlers of the Bot
	AfterInit   func(s *discordgo.Session) // AfterInit is called after the initialization process of the Bot
	Version     *Version
	Innovations []*Innovation
	Name        string
	Intents     discord.Intent
}

// Status contains all required information for updating the status
type Status struct {
	Type    StatusType // Type of the Status (use GameStatus or WatchStatus or StreamingStatus or ListeningStatus)
	Content string     // Content of the Status
	Url     string     // Url of the StreamingStatus
}

// Start the Bot (blocking instruction)
func (b *Bot) Start() {
	dg := discordgo.New("Bot " + b.Token) // New connection to the discord API with bot token

	dg.Identify.Intents = b.Intents
	b.Logger = dg

	if Debug {
		dg.ChangeLevel(logger.LevelDebug)
	}

	for _, handler := range b.handlers {
		dg.EventManager().AddHandler(handler)
	}

	err := dg.Open() // Starts the bot
	if err != nil {
		dg.LogError(err, "starting bot")
		return
	}

	var wg sync.WaitGroup
	st := time.Now()
	// register commands
	wg.Add(1)
	go func() {
		b.updateCommands(dg)
		dg.LogInfo("Commands updated in %s", time.Since(st))
		wg.Done()
	}()
	b.setupCommandsHandlers(dg)

	if Debug {
		dg.EventManager().AddHandler(func(s bot.Session, i event.InteractionCreate) {
			dg.LogDebug("Interaction received")
			data, _ := json.Marshal(i)
			dg.LogDebug("%s", data)
		})
	}
	if b.AfterInit != nil {
		b.AfterInit(dg)
	}

	// wait until all setup goroutines are finished
	wg.Wait()
	delta := time.Since(st)
	to := dg.Client.Timeout.Seconds()
	// if the setup was faster than the http client timeout, wait
	if delta.Seconds() < to {
		time.Sleep(time.Duration(to-delta.Seconds()) * time.Second)
	}
	dg.LogInfo("Bot started as %s", dg.SessionState().User().Username)
	NewTimer(30*time.Second, func(stop chan<- interface{}) {
		if b.Status == nil {
			stop <- struct{}{}
			return
		}
		l := len(b.Status)
		r := rand.New(rand.NewPCG(uint64(time.Now().Unix()), uint64(l))).UintN(uint(l))
		s := b.Status[r]
		if s.Type == GameStatus {
			err = dg.BotAPI().UpdateGameStatus(0, s.Content)
		} else if s.Type == WatchStatus {
			err = dg.BotAPI().UpdateWatchStatus(0, s.Content)
		} else if s.Type == StreamingStatus {
			if s.Url == "" {
				err = ErrStatusUrlNotFound
			} else {
				err = dg.BotAPI().UpdateStreamingStatus(0, s.Content, s.Url)
			}
		} else if s.Type == ListeningStatus {
			err = dg.BotAPI().UpdateListeningStatus(s.Content)
		} else {
			err = ErrBadStatusType
		}
		if err != nil {
			dg.LogError(err, "updating status")
			err = nil
		}
	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.LogInfo("Stopping bot")

	err = dg.Close() // Bot Shutdown
	if err != nil {
		dg.LogError(err, "closing bot")
	}

	dg.LogInfo("Bot shut down")
}

func (b *Bot) AddHandler(handler any) {
	b.handlers = append(b.handlers, handler)
}

func (b *Bot) HandleModal(handler func(bot.Session, *event.InteractionCreate, *interaction.ModalSubmitData, *cmd.ResponseBuilder),
	id string) {
	b.AddHandler(func(s bot.Session, i *event.InteractionCreate) {
		if i.Type != types.InteractionModalSubmit {
			return
		}

		data := i.ModalSubmitData()
		if data.CustomID != id {
			return
		}
		handler(s, i, data, cmd.NewResponseBuilder(s, i))
	})
}

func (b *Bot) HandleMessageComponent(handler func(bot.Session, *event.InteractionCreate, *interaction.MessageComponentData, *cmd.ResponseBuilder),
	id string) {
	b.AddHandler(func(s bot.Session, i *event.InteractionCreate) {
		if i.Type != types.InteractionMessageComponent {
			return
		}

		data := i.MessageComponentData()
		if data.CustomID != id {
			return
		}
		handler(s, i, data, cmd.NewResponseBuilder(s, i))
	})
}
