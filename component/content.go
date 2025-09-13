package component

import (
	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/component"
)

type Section struct {
	components []SubSection
	accessory  Accessory
	id         int
}

func (s *Section) SetID(i int) Sub {
	//TODO implement me
	panic("implement me")
}

func (s *Section) Component() component.Component {
	cp := make([]component.Component, len(s.components))
	for i, sub := range s.components {
		cp[i] = sub.Component()
	}
	return &component.Section{
		ID:         s.id,
		Components: cp,
		Accessory:  s.accessory.Component(),
	}
}

func (s *Section) IsForModal() bool {
	return false
}

func (s *Section) SetAccessory(accessory Accessory) *Section {
	s.accessory = accessory
	return s
}

func (s *Section) Add(sub SubSection) *Section {
	if s.components == nil {
		s.components = []SubSection{}
	}
	s.components = append(s.components, sub)
	return s
}

func (s *Section) subContainer() {}

func (s *Section) isTopLevel() {}

func NewSection() *Section {
	return new(Section)
}

type TextDisplay struct {
	component.TextDisplay
}

func (t *TextDisplay) Component() component.Component {
	return t.TextDisplay
}

func (t *TextDisplay) IsForModal() bool {
	return false
}

func (t *TextDisplay) SetID(i int) Sub {
	panic("Missing ID in discordgo.TextDisplay. gokord cannot fix this")
}

func (t *TextDisplay) SetContent(s string) *TextDisplay {
	t.Content = s
	return t
}

func (t *TextDisplay) subContainer() {}

func (t *TextDisplay) isTopLevel() {}

func (t *TextDisplay) subSection() {}

func NewTextDisplay(content string) *TextDisplay {
	return new(TextDisplay).SetContent(content)
}

type Thumbnail struct {
	component.Thumbnail
}

func (t *Thumbnail) Component() component.Component {
	return t.Thumbnail
}

func (t *Thumbnail) IsForModal() bool {
	return false
}

func (t *Thumbnail) SetID(i int) Sub {
	t.ID = i
	return t
}

func (t *Thumbnail) IsSpoiler() *Thumbnail {
	t.Spoiler = true
	return t
}

func (t *Thumbnail) SetDescription(s string) *Thumbnail {
	t.Description = &s
	return t
}

// SetMedia takes an URL
func (t *Thumbnail) SetMedia(s string) *Thumbnail {
	t.Media = component.UnfurledMediaItem{URL: s}
	return t
}

func (t *Thumbnail) accessory() {}

// NewThumbnail takes a URL as a media
func NewThumbnail(media string) *Thumbnail {
	return new(Thumbnail).SetMedia(media)
}

type MediaGallery struct {
	component.MediaGallery
}

func (m *MediaGallery) Component() component.Component {
	return m.MediaGallery
}

func (m *MediaGallery) IsForModal() bool {
	return false
}

func (m *MediaGallery) SetID(i int) Sub {
	m.ID = i
	return m
}

func (m *MediaGallery) Add(url string, description string, spoiler bool) *MediaGallery {
	if m.Items == nil {
		m.Items = []component.MediaGalleryItem{}
	}
	item := component.MediaGalleryItem{
		Media:       component.UnfurledMediaItem{URL: url},
		Description: &description,
		Spoiler:     spoiler,
	}
	if len(description) == 0 {
		item.Description = nil
	}
	m.Items = append(m.Items, item)
	return m
}

func (m *MediaGallery) subContainer() {}

func (m *MediaGallery) isTopLevel() {}

func NewMediaGallery() *MediaGallery {
	return new(MediaGallery)
}

type File struct {
	component.File
}

func (f *File) Component() component.Component {
	return f.File
}

func (f *File) IsForModal() bool {
	return false
}

func (f *File) SetID(i int) Sub {
	f.ID = i
	return f
}

func (f *File) IsSpoiler() *File {
	f.Spoiler = true
	return f
}

// SetFile takes an URL
func (f *File) SetFile(s string) *File {
	f.File.File = component.UnfurledMediaItem{URL: s}
	return f
}

func (f *File) subContainer() {}

func (f *File) isTopLevel() {}

// NewFile takes a URL as a media
func NewFile(media string) *File {
	return new(File).SetFile(media)
}

type Separator struct {
	component.Separator
}

func (s *Separator) Component() component.Component {
	return s.Separator
}

func (s *Separator) IsForModal() bool {
	return false
}

func (s *Separator) SetID(i int) Sub {
	s.ID = i
	return s
}

func (s *Separator) IsNotDivider() *Separator {
	b := false
	s.Divider = &b
	return s
}

func (s *Separator) SetSpacing(sp component.SeparatorSpacingSize) *Separator {
	s.Spacing = &sp
	return s
}

func (s *Separator) subContainer() {}

func (s *Separator) isTopLevel() {}

func NewSeparator() *Separator {
	s := new(Separator)
	b := true
	s.Divider = &b
	return s.SetSpacing(component.SeparatorSpacingSizeSmall)
}

type Container struct {
	components  []SubContainer
	id          int
	accentColor *int
	spoiler     bool
}

func (c *Container) Component() component.Component {
	cp := make([]component.Component, len(c.components))
	for i, c := range c.components {
		cp[i] = c.Component()
	}
	return &component.Container{
		Components:  cp,
		ID:          c.id,
		AccentColor: c.accentColor,
		Spoiler:     c.spoiler,
	}
}

func (c *Container) IsForModal() bool {
	return false
}

func (c *Container) SetID(i int) Sub {
	c.id = i
	return c
}

func (c *Container) IsSpoiler() *Container {
	c.spoiler = true
	return c
}

func (c *Container) SetAccentColor(i int) *Container {
	c.accentColor = &i
	return c
}

func (c *Container) Add(s SubContainer) *Container {
	if c.components == nil {
		c.components = []SubContainer{}
	}
	c.components = append(c.components, s)
	return c
}

func (c *Container) isTopLevel() {}

func NewContainer() *Container {
	return new(Container)
}
