package gokord

import (
	"strconv"
	"time"

	discordgo "github.com/nyttikord/gokord"
	"github.com/nyttikord/gokord/discord/types"
	"github.com/nyttikord/gokord/user"
)

// FetchGuildUser returns the list of member in a guild
func FetchGuildUser(s *discordgo.Session, guildID string) []*user.Member {
	member, err := s.GuildAPI().Members(guildID, "", 1000)
	if err != nil {
		s.LogError(err, "fetching guild users")
	}
	return member
}

func GetTimestampFromId(id string) (time.Time, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return time.UnixMilli(0), err
	}

	// https://discord.com/developers/docs/reference#snowflakes-snowflake-id-format-structure-left-to-right
	timestamp := (idInt >> 22) + 1420070400000

	return time.UnixMilli(timestamp), nil
}

// ComesFromDM returns true if a message comes from a DM channel
func ComesFromDM(s *discordgo.Session, id string) (bool, error) {
	channel, err := s.ChannelAPI().State.Channel(id)
	if err != nil {
		if channel, err = s.ChannelAPI().Channel(id); err != nil {
			return false, err
		}
	}

	return channel.Type == types.ChannelDM, nil
}
