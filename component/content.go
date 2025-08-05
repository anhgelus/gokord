package component

import "github.com/bwmarrin/discordgo"

type Section struct {
	components []*discordgo.TextDisplay
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

func (s *Section) AddComponent(c discordgo.MessageComponent) *Section {
	if s.components == nil {
		s.components = make([]*discordgo.TextDisplay, len(s.components))
	}
	s.components = append(s.components, c)
	return s
}
