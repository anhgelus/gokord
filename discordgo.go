package gokord

import (
	"github.com/anhgelus/gokord/logger"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

// FetchGuildUser returns the list of member in a guild
func FetchGuildUser(s *discordgo.Session, guildID string) []*discordgo.Member {
	member, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		logger.Alert("discordgo.go - Failed to fetch guild users", err.Error())
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
	channel, err := s.State.Channel(id)
	if err != nil {
		if channel, err = s.Channel(id); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}
