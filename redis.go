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

// RedisBase is an interface helping use of redis to store/cache data
type RedisBase interface {
	GenKey(key string) string // GenKey generates the key to use
}

// RedisUser is the default implementation of RedisBase for a Discord User
type RedisUser struct {
	RedisBase
	DiscordID string
	GuildID   string
}

var (
	ErrNilClient = errors.New("redis.NewClient is nil")
)

func (p *RedisUser) GenKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", p.GuildID, p.DiscordID, key)
}

// Connect to Redis with the given RedisCredentials
func (rc *RedisCredentials) Connect() (*redis.Client, error) {
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
