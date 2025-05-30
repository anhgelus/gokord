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
	files         []*discordgo.File
	//
	interaction *discordgo.InteractionCreate
	session     *discordgo.Session
}

func NewResponseBuilder(s *discordgo.Session, i *discordgo.InteractionCreate) *ResponseBuilder {
	return &ResponseBuilder{
		interaction: i,
		session:     s,
	}
}

// Send the response
func (res *ResponseBuilder) Send() error {
	if res.edit {
		_, err := res.session.InteractionResponseEdit(res.interaction.Interaction, &discordgo.WebhookEdit{
			Content: &res.content,
			Embeds:  &res.messageEmbeds,
			Files:   res.files,
		})
		return err
	}

	r := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: res.content,
			Embeds:  res.messageEmbeds,
			Files:   res.files,
		},
	}
	if res.deferred {
		r.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
	}
	if res.ephemeral {
		r.Data.Flags = discordgo.MessageFlagsEphemeral
	}

	if err := res.session.InteractionRespond(res.interaction.Interaction, r); err != nil {
		return err
	}

	if res.deferred {
		res.IsEdit()
	}
	return nil
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
	res.NotEdit()
	res.deferred = true
	return res
}

func (res *ResponseBuilder) NotDeferred() *ResponseBuilder {
	res.deferred = false
	return res
}

func (res *ResponseBuilder) IsEdit() *ResponseBuilder {
	res.NotDeferred()
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

func (res *ResponseBuilder) SetEmbeds(e []*discordgo.MessageEmbed) *ResponseBuilder {
	t := time.Now()
	footer := &discordgo.MessageEmbedFooter{
		Text:    "by " + Author,
		IconURL: res.session.State.User.AvatarURL(""),
	}
	author := &discordgo.MessageEmbedAuthor{
		Name: res.session.State.User.Username,
	}
	for _, em := range e {
		em.Footer = footer
		em.Timestamp = t.Format(time.RFC3339)
		em.Author = author
	}
	res.messageEmbeds = e
	return res
}

func (res *ResponseBuilder) SetFiles(f []*discordgo.File) *ResponseBuilder {
	res.files = f
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

func GenerateOptionMapForSubcommand(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options[0].Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
