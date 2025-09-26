package gokord

import (
	"errors"

	"github.com/anhgelus/gokord/cmd"
	"github.com/nyttikord/gokord/bot"
	"github.com/nyttikord/gokord/event"
)

var (
	ErrSubsAreNil         = errors.New("subs are nil in general handler")
	ErrSubCommandNotFound = errors.New("sub command not found")
)

// generalHandler used for subcommand
func (b *Bot) generalHandler(s bot.Session, i *event.InteractionCreate, _ cmd.OptionMap, resp *cmd.ResponseBuilder) {
	data := i.CommandData()
	sendWarn := func(msg string, msgSend string) {
		s.LogError(ErrSubCommandNotFound, "%s: %s", data.Name, msg)
		err := resp.IsEphemeral().SetMessage(msgSend).Send()
		if err != nil {
			s.LogError(err, "sending error")
		}
	}
	if len(data.Options) == 0 {
		sendWarn("len(data.Options) == 0", "No subcommand identified (may be a bug)")
		return
	}
	subInfo := data.Options[0]
	if subInfo == nil {
		sendWarn("subInfo == nil", "No subcommand identified")
		return
	}
	var c cmd.CommandBuilder
	for _, cb := range b.Commands {
		if cb.GetName() == data.Name {
			c = cb
		}
	}
	if c == nil {
		sendWarn("cmd == nil", "Command not found")
		return
	}
	if c.GetSubs() == nil {
		s.LogError(ErrSubsAreNil, "subs are nil in general handler")
		err := resp.IsEphemeral().SetMessage("Internal error, please report it").Send()
		if err != nil {
			s.LogError(err, "sending error")
		}
		return
	}
	for _, sub := range c.GetSubs() {
		if subInfo.Name == sub.GetName() {
			sub.GetHandler()(s, i, cmd.GenerateOptionMapForSubcommand(i), resp)
			return
		}
	}
	sendWarn("Subcommand not found", "Subcommand not found")
}
