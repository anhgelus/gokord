package component

import "github.com/bwmarrin/discordgo"

type GeneralContainer interface {
	Add(Sub) GeneralContainer
	Components() []discordgo.MessageComponent
	ForModal() GeneralContainer
}

type Sub interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	CanBeInContainer() bool
	SetID(int) Sub
}

type Interactive interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	CanBeInContainer() bool // must be false
	SetID(int) Sub
	SetCustomID(string) Interactive
}

type Accessory interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	CanBeInContainer() bool // must be false
	SetID(int) Sub
	accessory() // does nothing
}

type SubContainer interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	CanBeInContainer() bool // must be false
	SetID(int) Sub
	subContainer() // does nothing
}

type containerBuilder struct {
	subs  []Sub
	modal bool
}

func (b *containerBuilder) Add(sub Sub) GeneralContainer {
	if sub.CanBeInContainer() {
		panic("Sub component cannot be directly added in container")
	}
	if b.modal != sub.IsForModal() {
		if b.modal {
			panic("Sub component cannot be added for a modal component")
		}
		panic("Sub component cannot be added for a message component")
	}
	b.subs = append(b.subs, sub)
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
	return &containerBuilder{}
}
