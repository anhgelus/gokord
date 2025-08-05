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
	*discordgo.Button
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

type SelectOption struct {
	discordgo.SelectMenuOption
}

func (s *SelectOption) SetDescription(desc string) *SelectOption {
	s.Description = desc
	return s
}

func (s *SelectOption) SetEmoji(e *discordgo.ComponentEmoji) *SelectOption {
	s.Emoji = e
	return s
}

func (s *SelectOption) IsDefault() *SelectOption {
	s.Default = true
	return s
}

func NewSelectOption(label string, value string) *SelectOption {
	return &SelectOption{
		SelectMenuOption: discordgo.SelectMenuOption{
			Label: label,
			Value: value,
		},
	}
}

type StringSelect struct {
	*discordgo.SelectMenu
}

func (s *StringSelect) Component() discordgo.MessageComponent {
	s.MenuType = discordgo.StringSelectMenu
	return s.SelectMenu
}

func (s *StringSelect) IsForModal() bool {
	return false
}

func (s *StringSelect) CanBeInContainer() bool {
	return false
}

func (s *StringSelect) SetCustomID(id string) Interactive {
	s.CustomID = id
	return s
}

func (s *StringSelect) SetID(i int) Interactive {
	s.ID = i
	return s
}

func (s *StringSelect) IsDisabled() *StringSelect {
	s.Disabled = true
	return s
}

func (s *StringSelect) SetMinValues(i int) *StringSelect {
	s.MinValues = &i
	return s
}

func (s *StringSelect) SetMaxValues(i int) *StringSelect {
	s.MaxValues = i
	return s
}

func (s *StringSelect) SetPlaceholder(placeholder string) *StringSelect {
	s.Placeholder = placeholder
	return s
}

func (s *StringSelect) AddOption(opt *SelectOption) *StringSelect {
	if s.Options == nil {
		s.Options = []discordgo.SelectMenuOption{}
	}
	s.Options = append(s.Options, opt.SelectMenuOption)
	return s
}

type TextInput struct {
	*discordgo.TextInput
}

func (t *TextInput) Component() discordgo.MessageComponent {
	return t.TextInput
}

func (t *TextInput) IsForModal() bool {
	return true
}

func (t *TextInput) CanBeInContainer() bool {
	return false
}

func (t *TextInput) SetCustomID(s string) Interactive {
	t.CustomID = s
	return t
}

func (t *TextInput) SetID(i int) Interactive {
	t.ID = i
	return t
}

func (t *TextInput) SetStyle(s discordgo.TextInputStyle) *TextInput {
	t.Style = s
	return t
}

func (t *TextInput) SetMinLength(i int) *TextInput {
	t.MinLength = i
	return t
}

func (t *TextInput) SetMaxLength(i int) *TextInput {
	t.MaxLength = i
	return t
}

func (t *TextInput) SetLabel(s string) *TextInput {
	t.Label = s
	return t
}

func (t *TextInput) IsRequired() *TextInput {
	t.Required = true
	return t
}

func (t *TextInput) SetPlaceholder(placeholder string) *TextInput {
	t.Placeholder = placeholder
	return t
}

func (t *TextInput) SetValue(v string) *TextInput {
	t.Value = v
	return t
}

type UserSelect struct {
	discordgo.SelectMenu
}

func (u *UserSelect) Component() discordgo.MessageComponent {
	return u.SelectMenu
}

func (u *UserSelect) IsForModal() bool {
	return false
}

func (u *UserSelect) CanBeInContainer() bool {
	return false
}

func (u *UserSelect) SetCustomID(s string) Interactive {
	u.CustomID = s
	return u
}

func (u *UserSelect) SetID(i int) Interactive {
	u.ID = i
	return u
}

func (u *UserSelect) IsDisabled() *UserSelect {
	u.Disabled = true
	return u
}

func (u *UserSelect) SetMinValues(i int) *UserSelect {
	u.MinValues = &i
	return u
}

func (u *UserSelect) SetMaxValues(i int) *UserSelect {
	u.MaxValues = i
	return u
}

func (u *UserSelect) SetPlaceholder(placeholder string) *UserSelect {
	u.Placeholder = placeholder
	return u
}

// discordgo does not seem to support this
//func (u *UserSelect) AddOption(opt *SelectOption) *UserSelect {
//	if u.Options == nil {
//		u.Options = []discordgo.SelectMenuOption{}
//	}
//	u.Options = append(u.Options, opt.SelectMenuOption)
//	return u
//}
