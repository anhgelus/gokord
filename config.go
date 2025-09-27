package gokord

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/anhgelus/gokord/cmd"
	"github.com/pelletier/go-toml/v2"
	"gorm.io/gorm"
)

var (
	// BaseCfg is the main BaseConfig used by the bot
	BaseCfg BaseConfig

	// UseRedis is true if the bot will use redis
	UseRedis = true

	ErrImpossibleToConnectDB         = errors.New("impossible to connect to the database")
	ErrImpossibleToConnectRedis      = errors.New("impossible to connect to redis")
	ErrMigratingGokordInternalModels = errors.New("error while migrating internal models")
)

const (
	ConfigFolder = "config"
)

type BaseConfig interface {
	// IsDebug returns true if the bot is in debug mode
	IsDebug() bool
	// GetAuthor returns the author (or the owner) of the bot
	GetAuthor() string
	// GetRedisCredentials returns the RedisCredentials used by the bot.
	//
	// Must return nil if gokord.UseRedis is false
	GetRedisCredentials() *RedisCredentials
	// GetSQLCredentials returns the SQLCredentials used by the bot
	GetSQLCredentials() SQLCredentials
	// SetDefaultValues set all values of this config to their default ones. THIS IS A DESTRUCTIVE OPERATION!
	SetDefaultValues()
	// Marshal the config, must use toml.Marshal
	Marshal() ([]byte, error)
	// Unmarshal the config, must use toml.Unmarshal
	Unmarshal([]byte) error
}

type SQLCredentials interface {
	// SetDefaultValues set all values of these credentials to their default ones.
	// THIS IS A DESTRUCTIVE OPERATION!
	SetDefaultValues()
	// Connect to the database, must use gorm.Open
	// (see https://gorm.io/docs/connecting_to_the_database.html)
	Connect() (*gorm.DB, error)
}

type RedisCredentials struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// SetDefaultValues set all values of these credentials to their default ones.
// THIS IS A DESTRUCTIVE OPERATION!
func (rc *RedisCredentials) SetDefaultValues() {
	rc.Address = "localhost:6379"
	rc.Password = "password"
	rc.DB = 0
}

// ConfigInfo has all required information to get a config
type ConfigInfo struct {
	Cfg           interface{} // Cfg is a pointer to the struct
	Name          string      // Name of the config
	DefaultValues func()      // DefaultValues is called to set up the default values of the config
}

func setupBaseConfig() error {
	return LoadConfig(&BaseCfg, "config", BaseCfg.SetDefaultValues, func(_ interface{}) ([]byte, error) {
		return BaseCfg.Marshal()
	}, func(data []byte, _ interface{}) error {
		return BaseCfg.Unmarshal(data)
	})
}

// LoadConfig a config (already called on start)
func LoadConfig(cfg interface{}, name string, defaultValues func(), marshal func(interface{}) ([]byte, error), unmarshal func([]byte, interface{}) error) error {
	path := fmt.Sprintf("%s/%s.toml", ConfigFolder, name)
	err := os.Mkdir(ConfigFolder, 0666)
	if err != nil && !os.IsExist(err) {
		return err
	}
	c, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		// TEMP: better when the two gokord will be merged
		slog.Warn("File not found, creating a new one.")
		defaultValues()
		c, err = marshal(cfg)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, c, 0666)
		if err != nil {
			return err
		}
		return nil
	}
	return unmarshal(c, cfg)
}

// SetupConfigs with the given configs (+ base config which is available at BaseCfg)
//
// customBaseConfig is the new type of BaseCfg
func SetupConfigs(customBaseConfig BaseConfig, cfgInfo []*ConfigInfo) error {
	var err error
	BaseCfg = customBaseConfig
	err = setupBaseConfig()
	if err != nil {
		return err
	}

	Debug = BaseCfg.IsDebug()
	cmd.Author = BaseCfg.GetAuthor()
	if Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	for _, cfg := range cfgInfo {
		err = LoadConfig(cfg.Cfg, cfg.Name, cfg.DefaultValues, toml.Marshal, toml.Unmarshal)
		if err != nil {
			return err
		}
	}

	DB, err = BaseCfg.GetSQLCredentials().Connect()
	if err != nil {
		return errors.Join(ErrImpossibleToConnectDB, err)
	}

	err = DB.AutoMigrate(&BotData{})
	if err != nil {
		return errors.Join(ErrMigratingGokordInternalModels, err)
	}

	if !UseRedis {
		return nil
	}
	c, err := BaseCfg.GetRedisCredentials().Connect()
	if err != nil {
		return errors.Join(ErrImpossibleToConnectRedis, err)
	}
	_ = c.Close()
	return nil
}
