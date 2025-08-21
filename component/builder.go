package component

import "github.com/bwmarrin/discordgo"

type GeneralContainer interface {
	Add(TopLevel) GeneralContainer
	Components() []discordgo.MessageComponent
	ForModal() GeneralContainer
}

type Sub interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	SetID(int) Sub
}

type TopLevel interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	SetID(int) Sub
	isTopLevel()
}

type Interactive interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	SetID(int) Sub
	SetCustomID(string) Interactive
}

type Accessory interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	SetID(int) Sub
	accessory() // does nothing
}

type SubContainer interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	SetID(int) Sub
	subContainer() // does nothing
}

type SubSection interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	SetID(int) Sub
	subSection() // does nothing
}

type containerBuilder struct {
	subs  []TopLevel
	modal bool
}

func (b *containerBuilder) Add(t TopLevel) GeneralContainer {
	if b.modal != t.IsForModal() {
		if b.modal {
			panic("Top level component cannot be added for a modal component")
		}
		panic("Top level component cannot be added for a message component")
	}
	b.subs = append(b.subs, t)
	return b
}

func (b *containerBuilder) Components() []discordgo.MessageComponent {
	cp := make([]discordgo.MessageComponent, len(b.subs))
	for i, sub := range b.subs {
		cp[i] = sub.Component()
	}
	return cp
}

func (b *containerBuilder) ForModal() GeneralContainer {
	if len(b.subs) != 0 {
		panic("Cannot set for modal if subs are not empty")
	}
	b.modal = true
	return b
}

func New() GeneralContainer {
	return new(containerBuilder)
}
