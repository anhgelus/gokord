package gokord

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// Credentials of redis
var Credentials RedisCredentials

// Ctx background
var Ctx = context.Background()

var client *redis.Client

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
	if client != nil {
		c := redis.NewClient(&redis.Options{
			Addr:     rc.Address,
			Password: rc.Password,
			DB:       rc.DB,
		})
		if c == nil {
			return nil, ErrNilClient
		}
		err := client.Ping(context.Background()).Err()
		if err != nil {
			return nil, err
		}
		client = c
	}
	s := client.Ping(context.Background())
	var err error
	if s != nil {
		err = s.Err()
	}
	return client, err
}
