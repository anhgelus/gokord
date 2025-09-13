package gokord

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/anhgelus/gokord/cmd"
	"github.com/anhgelus/gokord/logger"
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/discord"
	"github.com/nyttikord/gokord/discord/types"
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

	err := dg.Open() // Starts the bot
	if err != nil {
		logger.Alert("bot.go - Start", err.Error())
		return
	}

	var wg sync.WaitGroup
	st1 := time.Now()
	// register commands
	wg.Add(1)
	go func() {
		b.updateCommands(dg)
		logger.Success("Commands updated")
		wg.Done()
	}()
	b.setupCommandsHandlers(dg)

	for _, handler := range b.handlers {
		dg.AddHandler(handler)
	}
	if Debug {
		dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			logger.Debug("Interaction received")
			data, _ := json.Marshal(i)
			logger.Debug(string(data))
		})
	}
	if b.AfterInit != nil {
		b.AfterInit(dg)
	}

	// wait until all setup goroutines are finished
	wg.Wait()
	st2 := time.Now()
	delta := float64(st2.Unix() - st1.Unix())
	to := dg.Client.Timeout.Seconds()
	// if the setup was faster than the http client timeout, wait
	if delta < dg.Client.Timeout.Seconds() {
		time.Sleep(time.Duration(to-delta) * time.Second)
	}
	logger.Success(fmt.Sprintf("Bot started as %s", dg.State.User.Username))
	NewTimer(30*time.Second, func(stop chan<- interface{}) {
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
			logger.Alert("bot.go - Update status", err.Error())
			err = nil
		}
	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	logger.Success("Bot shutting down")

	if Debug {
		logger.Debug("Unregistering local commands")
		b.unregisterGuildCommands(dg)
	}

	err = dg.Close() // Bot Shutdown
	if err != nil {
		logger.Alert("bot.go - Shutdown", err.Error())
	}

	logger.Success("Bot shut down")
}

func (b *Bot) AddHandler(handler any) {
	b.handlers = append(b.handlers, handler)
}

func (b *Bot) HandleModal(handler func(*discordgo.Session, *discordgo.InteractionCreate, interaction.ModalSubmitInteractionData, *cmd.ResponseBuilder),
	id string) {
	b.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

func (b *Bot) HandleMessageComponent(handler func(*discordgo.Session, *discordgo.InteractionCreate, interaction.MessageComponentInteractionData, *cmd.ResponseBuilder),
	id string) {
	b.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
