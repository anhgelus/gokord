package utils

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

var Author string

// ResponseBuilder helps to response to slash commands
type ResponseBuilder struct {
	content       string
	ephemeral     bool
	deferred      bool
	edit          bool
	messageEmbeds []*discordgo.MessageEmbed
	I             *discordgo.InteractionCreate
	C             *discordgo.Session
}

// Send the response
func (res *ResponseBuilder) Send() error {
	if res.edit {
		_, err := res.C.InteractionResponseEdit(res.I.Interaction, &discordgo.WebhookEdit{
			Content: &res.content,
			Embeds:  &res.messageEmbeds,
		})
		return err
	}
	r := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res.content,
			Embeds:  res.messageEmbeds,
		},
	}
	if res.deferred {
		r.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
	}
	if res.ephemeral {
		r.Data.Flags = discordgo.MessageFlagsEphemeral
	}
	return res.C.InteractionRespond(res.I.Interaction, r)
}

func (res *ResponseBuilder) Interaction(i *discordgo.InteractionCreate) *ResponseBuilder {
	res.I = i
	return res
}

func (res *ResponseBuilder) Client(c *discordgo.Session) *ResponseBuilder {
	res.C = c
	return res
}

func (res *ResponseBuilder) IsEphemeral() *ResponseBuilder {
	res.ephemeral = true
	return res
}

func (res *ResponseBuilder) NotEphemeral() *ResponseBuilder {
	res.ephemeral = false
	return res
}

func (res *ResponseBuilder) IsDeferred() *ResponseBuilder {
	res.deferred = true
	return res
}

func (res *ResponseBuilder) NotDeferred() *ResponseBuilder {
	res.deferred = false
	return res
}

func (res *ResponseBuilder) IsEdit() *ResponseBuilder {
	res.edit = true
	return res
}

func (res *ResponseBuilder) NotEdit() *ResponseBuilder {
	res.edit = false
	return res
}

func (res *ResponseBuilder) Message(s string) *ResponseBuilder {
	res.content = s
	return res
}

func (res *ResponseBuilder) Embeds(e []*discordgo.MessageEmbed) *ResponseBuilder {
	t := time.Now()
	footer := &discordgo.MessageEmbedFooter{
		Text:    "by " + Author,
		IconURL: res.C.State.User.AvatarURL(""),
	}
	author := &discordgo.MessageEmbedAuthor{
		Name: res.C.State.User.Username,
	}
	for _, em := range e {
		em.Footer = footer
		em.Timestamp = t.Format(time.RFC3339)
		em.Author = author
	}
	res.messageEmbeds = e
	return res
}

func GenerateOptionMap(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
