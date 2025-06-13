package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

var Author string

// ResponseBuilder helps to response to slash commands
type ResponseBuilder struct {
	content    string
	ephemeral  bool
	deferred   bool
	edit       bool
	modal      bool
	components []discordgo.MessageComponent
	embeds     []*discordgo.MessageEmbed
	files      []*discordgo.File
	title      string
	customID   string
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

func formatInteractionResponse(r interface{}) string {
	b, e := json.MarshalIndent(r, "", "  ")
	if e != nil {
		panic(e)
	}
	return string(b)
}

// Send the response
func (res *ResponseBuilder) Send() error {
	if res.edit {
		wb := &discordgo.WebhookEdit{
			Content:    &res.content,
			Components: &res.components,
			Embeds:     &res.embeds,
			Files:      res.files,
		}
		_, err := res.session.InteractionResponseEdit(res.interaction.Interaction, wb)
		if err != nil {
			fmt.Println(formatInteractionResponse(wb))
		}
		return err
	}

	r := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    res.content,
			Components: res.components,
			Embeds:     res.embeds,
			Files:      res.files,
			CustomID:   res.customID,
			Title:      res.title,
		},
	}
	if res.deferred {
		r.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
	}
	if res.ephemeral {
		r.Data.Flags = discordgo.MessageFlagsEphemeral
	}
	if res.modal {
		r.Type = discordgo.InteractionResponseModal
	}

	if err := res.session.InteractionRespond(res.interaction.Interaction, r); err != nil {
		fmt.Println(formatInteractionResponse(r))
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
	res.NotModal()
	res.deferred = true
	return res
}

func (res *ResponseBuilder) NotDeferred() *ResponseBuilder {
	res.deferred = false
	return res
}

func (res *ResponseBuilder) IsEdit() *ResponseBuilder {
	res.NotDeferred()
	res.NotModal()
	res.edit = true
	return res
}

func (res *ResponseBuilder) NotEdit() *ResponseBuilder {
	res.edit = false
	return res
}

func (res *ResponseBuilder) IsModal() *ResponseBuilder {
	res.NotDeferred()
	res.NotEdit()
	res.NotEphemeral()
	res.modal = true
	return res
}

func (res *ResponseBuilder) NotModal() *ResponseBuilder {
	res.modal = false
	return res
}

func (res *ResponseBuilder) SetMessage(s string) *ResponseBuilder {
	res.content = s
	return res
}

func (res *ResponseBuilder) SetTitle(s string) *ResponseBuilder {
	res.title = s
	return res
}

func (res *ResponseBuilder) SetCustomID(s string) *ResponseBuilder {
	res.customID = s
	return res
}

func (res *ResponseBuilder) AddEmbed(e *discordgo.MessageEmbed) *ResponseBuilder {
	t := time.Now()
	e.Footer = &discordgo.MessageEmbedFooter{
		Text:    "by " + Author,
		IconURL: res.session.State.User.AvatarURL(""),
	}
	e.Timestamp = t.Format(time.RFC3339)
	e.Author = &discordgo.MessageEmbedAuthor{
		Name: res.session.State.User.Username,
	}
	if res.embeds == nil {
		res.embeds = []*discordgo.MessageEmbed{e}
	} else {
		res.embeds = append(res.embeds, e)
	}
	return res
}

func (res *ResponseBuilder) AddFile(f *discordgo.File) *ResponseBuilder {
	if res.files == nil {
		res.files = []*discordgo.File{f}
	} else {
		res.files = append(res.files, f)
	}
	return res
}

func (res *ResponseBuilder) AddComponent(c discordgo.MessageComponent) *ResponseBuilder {
	if res.components == nil {
		res.components = []discordgo.MessageComponent{c}
	} else {
		res.components = append(res.components, c)
	}
	return res
}

type OptionMap map[string]*discordgo.ApplicationCommandInteractionDataOption

func GenerateOptionMap(i *discordgo.InteractionCreate) OptionMap {
	options := i.ApplicationCommandData().Options
	optionMap := make(OptionMap, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func GenerateOptionMapForSubcommand(i *discordgo.InteractionCreate) OptionMap {
	options := i.ApplicationCommandData().Options[0].Options
	optionMap := make(OptionMap, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
