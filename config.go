package gokord

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/anhgelus/gokord/utils"
	"github.com/pelletier/go-toml/v2"
	"log/slog"
	"os"
)

var (
	//go:embed resources/config.toml
	DefaultBaseConfig string

	//go:embed resources/config_no_redis.toml
	DefaultBaseConfigNoRedis string

	// BaseCfg is the BaseConfig used by the bot
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

// BaseConfig is all basic configuration (debug, redis connection and database connection)
type BaseConfig struct {
	Debug    bool                `toml:"debug"`
	Author   string              `toml:"author"`
	Redis    RedisCredentials    `toml:"redis"`
	Database DatabaseCredentials `toml:"database"`
}

type RedisCredentials struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type DatabaseCredentials struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
	Port     int    `toml:"port"`
}

// ConfigInfo has all required information to get a config
type ConfigInfo struct {
	Cfg     any    // pointer to the struct
	Name    string // name of the config
	Default string // default content of the config
}

// getBaseConfig get the BaseConfig
func getBaseConfig(cfg any, defaultConfig string) error {
	return Get(cfg, defaultConfig, "config")
}

// Get a config (already called on start)
func Get(cfg any, defaultConfig string, name string) error {
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
		utils.SendAlert("config.go - Create file", "File not found, creating a new one.")
		c = []byte(defaultConfig)
		err = os.WriteFile(path, c, 0666)
		if err != nil {
			return err
		}
	}
	return toml.Unmarshal(c, cfg)
}

// SetupConfigs with the given configs (+ base config which is available at BaseCfg)
func SetupConfigs(cfgInfo []*ConfigInfo) error {
	var err error
	if UseRedis {
		err = getBaseConfig(&BaseCfg, DefaultBaseConfig)
	} else {
		err = getBaseConfig(&BaseCfg, DefaultBaseConfigNoRedis)
	}
	if err != nil {
		return err
	}

	Debug = BaseCfg.Debug
	utils.Author = BaseCfg.Author
	if Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	for _, cfg := range cfgInfo {
		err = Get(cfg.Cfg, cfg.Name, cfg.Default)
		if err != nil {
			return err
		}
	}

	DB, err = BaseCfg.Database.Connect()
	if err != nil {
		utils.SendAlert("config.go - connection to database", err.Error())
		return ErrImpossibleToConnectDB
	}

	err = DB.AutoMigrate(&BotData{})
	if err != nil {
		utils.SendAlert("config.go - migrating internal models", err.Error())
		return ErrMigratingGokordInternalModels
	}

	if !UseRedis {
		return nil
	}
	c, err := BaseCfg.Redis.Get()
	if err != nil {
		utils.SendAlert("config.go - connection to redis", err.Error())
		return ErrImpossibleToConnectRedis
	}
	_ = c.Close()
	return nil
}
