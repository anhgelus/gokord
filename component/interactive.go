package component

import (
	"github.com/bwmarrin/discordgo"
)

type ActionRow struct {
	subs  []Sub
	modal bool
}

func (a *ActionRow) Component() discordgo.MessageComponent {
	return discordgo.ActionsRow{
		Components: a.Components(),
	}
}

func (a *ActionRow) IsForModal() bool {
	return a.modal
}

func (a *ActionRow) CanBeInContainer() bool {
	return true
}

func (a *ActionRow) Add(sub Sub) {
	a.subs = append(a.subs, sub)
}

func (a *ActionRow) Components() []discordgo.MessageComponent {
	cp := make([]discordgo.MessageComponent, len(a.subs))
	for i, sub := range a.subs {
		cp[i] = sub.Component()
	}
	return cp
}

func (a *ActionRow) ForModal() {
	if len(a.subs) != 0 {
		panic("Cannot set for modal if subs are not empty")
	}
	a.modal = true
}

type Button struct {
	discordgo.Button
}

func (b *Button) Component() discordgo.MessageComponent {
	return b.Button
}

func (b *Button) IsForModal() bool {
	return false
}

func (b *Button) CanBeInContainer() bool {
	return false
}

func (b *Button) SetCustomID(s string) Interactive {
	b.CustomID = s
	return b
}

func (b *Button) SetID(i int) Interactive {
	b.ID = i
	return b
}

func (b *Button) SetLabel(l string) *Button {
	b.Label = l
	return b
}

func (b *Button) SetStyle(s discordgo.ButtonStyle) *Button {
	b.Style = s
	return b
}

func (b *Button) IsDisabled() *Button {
	b.Disabled = true
	return b
}

func (b *Button) SetEmoji(e *discordgo.ComponentEmoji) *Button {
	b.Emoji = e
	return b
}

func (b *Button) SetURL(url string) *Button {
	b.URL = url
	return b
}

func (b *Button) accessory() {}
