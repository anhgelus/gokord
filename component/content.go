package component

import "github.com/bwmarrin/discordgo"

type Section struct {
	components []Sub
	accessory  Accessory
	id         int
}

func (s *Section) SetID(i int) Sub {
	//TODO implement me
	panic("implement me")
}

func (s *Section) Component() discordgo.MessageComponent {
	cp := make([]discordgo.MessageComponent, len(s.components))
	for i, sub := range s.components {
		cp[i] = sub.Component()
	}
	return discordgo.Section{
		ID:         s.id,
		Components: cp,
		Accessory:  s.accessory.Component(),
	}
}

func (s *Section) IsForModal() bool {
	return false
}

func (s *Section) CanBeInContainer() bool {
	return true
}

func (s *Section) SetAccessory(accessory Accessory) *Section {
	s.accessory = accessory
	return s
}

func (s *Section) Add(sub Sub) *Section {
	if s.components == nil {
		s.components = make([]Sub, len(s.components))
	}
	s.components = append(s.components, sub)
	return s
}

func (s *Section) subContainer() {}

type TextDisplay struct {
	discordgo.TextDisplay
}

func (t *TextDisplay) Component() discordgo.MessageComponent {
	return t.TextDisplay
}

func (t *TextDisplay) IsForModal() bool {
	return false
}

func (t *TextDisplay) CanBeInContainer() bool {
	return true
}

func (t *TextDisplay) SetID(i int) Sub {
	panic("Missing ID in discordgo.TextDisplay. gokord cannot fix this")
}

func (t *TextDisplay) SetContent(s string) *TextDisplay {
	t.Content = s
	return t
}

func (s *TextDisplay) subContainer() {}

type Thumbnail struct {
	discordgo.Thumbnail
}

func (t *Thumbnail) Component() discordgo.MessageComponent {
	return t.Thumbnail
}

func (t *Thumbnail) IsForModal() bool {
	return false
}

func (t *Thumbnail) CanBeInContainer() bool {
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
	t.Media = discordgo.UnfurledMediaItem{URL: s}
	return t
}

func (t *Thumbnail) accessory() {}

type MediaGallery struct {
	discordgo.MediaGallery
}

func (m *MediaGallery) Component() discordgo.MessageComponent {
	return m.MediaGallery
}

func (m *MediaGallery) IsForModal() bool {
	return false
}

func (m *MediaGallery) CanBeInContainer() bool {
	return true
}

func (m *MediaGallery) SetID(i int) Sub {
	m.ID = i
	return m
}

func (m *MediaGallery) Add(url string, description string, spoiler bool) *MediaGallery {
	if m.Items == nil {
		m.Items = []discordgo.MediaGalleryItem{}
	}
	item := discordgo.MediaGalleryItem{
		Media:       discordgo.UnfurledMediaItem{URL: url},
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

type File struct {
	discordgo.FileComponent
}

func (f *File) Component() discordgo.MessageComponent {
	return f.FileComponent
}

func (f *File) IsForModal() bool {
	return false
}

func (f *File) CanBeInContainer() bool {
	return true
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
	f.File = discordgo.UnfurledMediaItem{URL: s}
	return f
}

func (f *File) subContainer() {}

type Container struct {
	components  []SubContainer
	id          int
	accentColor *int
	spoiler     bool
}

func (c *Container) Component() discordgo.MessageComponent {
	cp := make([]discordgo.MessageComponent, len(c.components))
	for i, c := range c.components {
		cp[i] = c.Component()
	}
	return discordgo.Container{
		Components:  cp,
		ID:          c.id,
		AccentColor: c.accentColor,
		Spoiler:     c.spoiler,
	}
}

func (c *Container) IsForModal() bool {
	return false
}

func (c *Container) CanBeInContainer() bool {
	return true
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
