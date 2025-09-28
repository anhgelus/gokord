package gokord

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anhgelus/gokord/cmd"
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/event"
	"github.com/nyttikord/gokord/interaction"
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
	Logger      *slog.Logger
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
	level := slog.LevelInfo
	if Debug {
		level = slog.LevelDebug
	}
	dg := discordgo.NewWithLogLevel("Bot "+b.Token, level) // New connection to the discord API with bot token

	dg.Identify.Intents = b.Intents
	b.Logger = dg.Logger()

	dg.EventManager().AddHandler(b.onReady)

	for _, handler := range b.handlers {
		dg.EventManager().AddHandler(handler)
	}

	err := dg.Open() // Starts the bot
	if err != nil {
		b.Logger.Error("starting bot", "error", err)
		return
	}

	// register commands
	go func() {
		st := time.Now()
		b.updateCommands(dg)
		b.Logger.Info("commands updated", "in", time.Since(st))
	}()
	b.setupCommandsHandlers(dg)

	if Debug {
		dg.EventManager().AddHandler(func(s bot.Session, i *event.InteractionCreate) {
			b.Logger.Debug("interaction received")
			data, _ := json.Marshal(i)
			b.Logger.Debug(string(data))
		})
	}
	if b.AfterInit != nil {
		b.AfterInit(dg)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	b.Logger.Info("stopping bot")

	err = dg.Close() // Bot Shutdown
	if err != nil {
		b.Logger.Error("closing bot", "error", err)
		b.Logger.Warn("force closing")
		dg.ForceClose()
	}

	b.Logger.Info("bot shut down")
}

func (b *Bot) onReady(s bot.Session, r *event.Ready) {
	b.Logger.Info("bot started", "as", s.SessionState().User().Username)
	NewTimer(30*time.Second, func(stop chan<- interface{}) {
		if b.Status == nil {
			stop <- struct{}{}
			return
		}
		l := len(b.Status)
		rnd := rand.New(rand.NewPCG(uint64(time.Now().Unix()), uint64(l))).UintN(uint(l))
		status := b.Status[rnd]
		var err error
		if status.Type == GameStatus {
			err = s.BotAPI().UpdateGameStatus(0, status.Content)
		} else if status.Type == WatchStatus {
			err = s.BotAPI().UpdateWatchStatus(0, status.Content)
		} else if status.Type == StreamingStatus {
			if status.Url == "" {
				err = ErrStatusUrlNotFound
			} else {
				err = s.BotAPI().UpdateStreamingStatus(0, status.Content, status.Url)
			}
		} else if status.Type == ListeningStatus {
			err = s.BotAPI().UpdateListeningStatus(status.Content)
		} else {
			err = ErrBadStatusType
		}
		if err != nil {
			b.Logger.Error("updating status", "error", err)
			err = nil
		}
	})
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
