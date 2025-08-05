package component

import "github.com/bwmarrin/discordgo"

type Container interface {
	Add(Sub)
	Components() []discordgo.MessageComponent
	ForModal()
}

type Sub interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	CanBeInContainer() bool
}

type Interactive interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	CanBeInContainer() bool // must be false
	interactive()           // does nothing
}

type Accessory interface {
	Component() discordgo.MessageComponent
	IsForModal() bool
	CanBeInContainer() bool // must be false
	accessory()             // does nothing
}

type containerBuilder struct {
	subs  []Sub
	modal bool
}

func (b *containerBuilder) Add(sub Sub) {
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
}

func (b *containerBuilder) Components() []discordgo.MessageComponent {
	cp := make([]discordgo.MessageComponent, len(b.subs))
	for i, sub := range b.subs {
		cp[i] = sub.Component()
	}
	return cp
}

func (b *containerBuilder) ForModal() {
	b.modal = true
}

func NewContainer() Container {
	return &containerBuilder{}
}
