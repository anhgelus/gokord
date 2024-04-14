package gokord

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	// Credentials of redis
	Credentials RedisCredentials
	// Ctx background
	Ctx = context.Background()
)

// UserBase is the minimum struct required to store a used in redis
type UserBase struct {
	DiscordID string
	GuildID   string
}

var (
	ErrGuildIDDiscordIDNotPresent = errors.New("guild_id or discord_id not informed")
	ErrNilClient                  = errors.New("redis.NewClient is nil")
)

func (p *UserBase) GenKey() string {
	return fmt.Sprintf("%s:%s", p.GuildID, p.DiscordID)
}

// Get the redis.Client with the given RedisCredentials
func (rc *RedisCredentials) Get() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     rc.Address,
		Password: rc.Password,
		DB:       rc.DB,
	})
	if client == nil {
		return nil, ErrNilClient
	}
	err := client.Ping(Ctx).Err()
	if err != nil {
		return nil, err
	}
	return client, err
}
