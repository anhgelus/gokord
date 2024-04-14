package gokord

import (
	"errors"
	"fmt"
	"github.com/anhgelus/gokord/utils"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Debug bool

	ErrBadStatusType     = errors.New("bad status type, please use the constant")
	ErrStatusUrlNotFound = errors.New("status url not found")
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
	Token    string        // Token of the Bot
	Status   []*Status     // Status of the Bot
	Commands []*Cmd        // Commands of the Bot
	Handlers []interface{} // Handlers of the Bot
}

// Status contains all required information for updating the status
type Status struct {
	Type    int    // Type of the Status (use GameStatus or WatchStatus or StreamingStatus or ListeningStatus)
	Content string // Content of the Status
	Url     string // Url of the StreamingStatus
}

// Start the Bot
func (b *Bot) Start() {
	dg, err := discordgo.New("Bot " + b.Token) // Define connection to discord API with bot token
	if err != nil {
		utils.SendAlert("bot.go - Token", err.Error())
	}

	err = dg.Open() // Bot start
	if err != nil {
		utils.SendAlert("bot.go - Start", err.Error())
	}
	go func() {
		time.Sleep(30 * time.Second)
		utils.SendSuccess(fmt.Sprintf("Bot started as %s", dg.State.User.Username))
		utils.NewTimer(30*time.Second, func(stop chan struct{}) {
			if b.Status == nil {
				stop <- struct{}{}
				return
			}
			rand.NewSource(time.Now().Unix())
			r := rand.Intn(len(b.Status))
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
	}()

	go func() {
		b.registerCommands(dg)
		utils.SendSuccess("Commands registered")
	}()
	b.setupCommandsHandlers(dg)

	for h := range b.Handlers {
		dg.AddHandler(h)
	}

	dg.Identify.Intents = discordgo.IntentsAll

	dg.StateEnabled = true

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	utils.SendDebug("Unregistering local commands")
	b.unregisterCommands(dg)

	err = dg.Close() // Bot Shutdown
	if err != nil {
		utils.SendAlert("bot.go - Shutdown", err.Error())
	}
}
