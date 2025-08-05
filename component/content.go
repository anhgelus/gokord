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

func (s *Section) AddComponent(sub Sub) *Section {
	if s.components == nil {
		s.components = make([]Sub, len(s.components))
	}
	s.components = append(s.components, sub)
	return s
}

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
