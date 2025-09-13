package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anhgelus/gokord/component"
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/channel"
	nc "github.com/nyttikord/gokord/component"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/interaction"
)

var Author string

// ResponseBuilder helps to response to slash commands
type ResponseBuilder struct {
	content    string
	ephemeral  bool
	deferred   bool
	edit       bool
	modal      bool
	components []nc.Component
	embeds     []*channel.MessageEmbed
	files      []*channel.File
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
		wb := &channel.WebhookEdit{
			Content:    &res.content,
			Components: &res.components,
			Embeds:     &res.embeds,
			Files:      res.files,
		}
		_, err := res.session.InteractionAPI().ResponseEdit(res.interaction.Interaction, wb)
		if err != nil {
			fmt.Println(formatInteractionResponse(wb))
		}
		return err
	}

	r := &interaction.InteractionResponse{
		Type: types.InteractionResponseChannelMessageWithSource,
		Data: &interaction.InteractionResponseData{
			Content:    res.content,
			Components: res.components,
			Embeds:     res.embeds,
			Files:      res.files,
			CustomID:   res.customID,
			Title:      res.title,
		},
	}
	if res.deferred {
		r.Type = types.InteractionResponseDeferredChannelMessageWithSource
	}
	if res.ephemeral {
		r.Data.Flags = channel.MessageFlagsEphemeral
	}
	if res.modal {
		r.Type = types.InteractionResponseModal
	}

	if err := res.session.InteractionAPI().Respond(res.interaction.Interaction, r); err != nil {
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

func (res *ResponseBuilder) AddEmbed(e *channel.MessageEmbed) *ResponseBuilder {
	t := time.Now()
	e.Footer = &channel.MessageEmbedFooter{
		Text:    "by " + Author,
		IconURL: res.session.State.User.AvatarURL(""),
	}
	e.Timestamp = t.Format(time.RFC3339)
	e.Author = &channel.MessageEmbedAuthor{
		Name: res.session.State.User.Username,
	}
	if res.embeds == nil {
		res.embeds = []*channel.MessageEmbed{e}
	} else {
		res.embeds = append(res.embeds, e)
	}
	return res
}

func (res *ResponseBuilder) AddFile(f *channel.File) *ResponseBuilder {
	if res.files == nil {
		res.files = []*channel.File{f}
	} else {
		res.files = append(res.files, f)
	}
	return res
}

func (res *ResponseBuilder) SetComponents(c component.GeneralContainer) *ResponseBuilder {
	res.components = c.Components()
	return res
}
