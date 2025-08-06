package component

import (
	"github.com/bwmarrin/discordgo"
)

type ActionRow struct {
	subs  []Sub
	modal bool
	id    int
}

func (a *ActionRow) SetID(i int) Sub {
	a.id = i
	return a
}

func (a *ActionRow) Component() discordgo.MessageComponent {
	cp := make([]discordgo.MessageComponent, len(a.subs))
	for i, sub := range a.subs {
		cp[i] = sub.Component()
	}
	return discordgo.ActionsRow{
		Components: cp,
	}
}

func (a *ActionRow) IsForModal() bool {
	return a.modal
}

func (a *ActionRow) CanBeInContainer() bool {
	return true
}

func (a *ActionRow) Add(sub Sub) *ActionRow {
	a.subs = append(a.subs, sub)
	return a
}

func (a *ActionRow) ForModal() {
	if len(a.subs) != 0 {
		panic("Cannot set for modal if subs are not empty")
	}
	a.modal = true
}

func (a *ActionRow) subContainer() {}

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

func (b *Button) SetID(i int) Sub {
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

func (s *StringSelect) SetID(i int) Sub {
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

func (t *TextInput) SetID(i int) Sub {
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
	u.MenuType = discordgo.UserSelectMenu
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

func (u *UserSelect) SetID(i int) Sub {
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

func (u *UserSelect) AddDefault(id string) *UserSelect {
	if u.DefaultValues == nil {
		u.DefaultValues = []discordgo.SelectMenuDefaultValue{}
	}
	u.DefaultValues = append(u.DefaultValues, discordgo.SelectMenuDefaultValue{ID: id, Type: discordgo.SelectMenuDefaultValueUser})
	return u
}

type RoleSelect struct {
	discordgo.SelectMenu
}

func (r *RoleSelect) Component() discordgo.MessageComponent {
	r.MenuType = discordgo.RoleSelectMenu
	return r.SelectMenu
}

func (r *RoleSelect) IsForModal() bool {
	return false
}

func (r *RoleSelect) CanBeInContainer() bool {
	return false
}

func (r *RoleSelect) SetCustomID(s string) Interactive {
	r.CustomID = s
	return r
}

func (r *RoleSelect) SetID(i int) Sub {
	r.ID = i
	return r
}

func (r *RoleSelect) IsDisabled() *RoleSelect {
	r.Disabled = true
	return r
}

func (r *RoleSelect) SetMinValues(i int) *RoleSelect {
	r.MinValues = &i
	return r
}

func (r *RoleSelect) SetMaxValues(i int) *RoleSelect {
	r.MaxValues = i
	return r
}

func (r *RoleSelect) SetPlaceholder(placeholder string) *RoleSelect {
	r.Placeholder = placeholder
	return r
}

func (r *RoleSelect) AddDefault(id string) *RoleSelect {
	if r.DefaultValues == nil {
		r.DefaultValues = []discordgo.SelectMenuDefaultValue{}
	}
	r.DefaultValues = append(r.DefaultValues, discordgo.SelectMenuDefaultValue{ID: id, Type: discordgo.SelectMenuDefaultValueRole})
	return r
}

type MentionableSelect struct {
	discordgo.SelectMenu
}

func (m *MentionableSelect) Component() discordgo.MessageComponent {
	m.MenuType = discordgo.MentionableSelectMenu
	return m.SelectMenu
}

func (m *MentionableSelect) IsForModal() bool {
	return false
}

func (m *MentionableSelect) CanBeInContainer() bool {
	return false
}

func (m *MentionableSelect) SetCustomID(s string) Interactive {
	m.CustomID = s
	return m
}

func (m *MentionableSelect) SetID(i int) Sub {
	m.ID = i
	return m
}

func (m *MentionableSelect) IsDisabled() *MentionableSelect {
	m.Disabled = true
	return m
}

func (m *MentionableSelect) SetMinValues(i int) *MentionableSelect {
	m.MinValues = &i
	return m
}

func (m *MentionableSelect) SetMaxValues(i int) *MentionableSelect {
	m.MaxValues = i
	return m
}

func (m *MentionableSelect) SetPlaceholder(placeholder string) *MentionableSelect {
	m.Placeholder = placeholder
	return m
}

func (m *MentionableSelect) AddDefault(id string, tp discordgo.SelectMenuDefaultValueType) *MentionableSelect {
	if m.DefaultValues == nil {
		m.DefaultValues = []discordgo.SelectMenuDefaultValue{}
	}
	m.DefaultValues = append(m.DefaultValues, discordgo.SelectMenuDefaultValue{ID: id, Type: tp})
	return m
}

type ChannelSelect struct {
	discordgo.SelectMenu
}

func (m *ChannelSelect) Component() discordgo.MessageComponent {
	m.MenuType = discordgo.ChannelSelectMenu
	return m.SelectMenu
}

func (m *ChannelSelect) IsForModal() bool {
	return false
}

func (m *ChannelSelect) CanBeInContainer() bool {
	return false
}

func (m *ChannelSelect) SetCustomID(s string) Interactive {
	m.CustomID = s
	return m
}

func (m *ChannelSelect) SetID(i int) Sub {
	m.ID = i
	return m
}

func (m *ChannelSelect) IsDisabled() *ChannelSelect {
	m.Disabled = true
	return m
}

func (m *ChannelSelect) SetMinValues(i int) *ChannelSelect {
	m.MinValues = &i
	return m
}

func (m *ChannelSelect) SetMaxValues(i int) *ChannelSelect {
	m.MaxValues = i
	return m
}

func (m *ChannelSelect) SetPlaceholder(placeholder string) *ChannelSelect {
	m.Placeholder = placeholder
	return m
}

func (m *ChannelSelect) AddDefault(id string) *ChannelSelect {
	if m.DefaultValues == nil {
		m.DefaultValues = []discordgo.SelectMenuDefaultValue{}
	}
	m.DefaultValues = append(m.DefaultValues, discordgo.SelectMenuDefaultValue{ID: id, Type: discordgo.SelectMenuDefaultValueChannel})
	return m
}

func (m *ChannelSelect) AddChannelType(tp discordgo.ChannelType) *ChannelSelect {
	if m.ChannelTypes == nil {
		m.ChannelTypes = []discordgo.ChannelType{}
	}
	m.ChannelTypes = append(m.ChannelTypes, tp)
	return m
}
