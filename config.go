package gokord

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/anhgelus/gokord/utils"
	"github.com/pelletier/go-toml/v2"
	"os"
)

var (
	//go:embed resources/config.toml
	DefaultBaseConfig string

	// BaseCfg is the BaseConfig used by the bot
	BaseCfg BaseConfig

	ErrImpossibleToConnectDB    = errors.New("impossible to connect to the database")
	ErrImpossibleToConnectRedis = errors.New("impossible to connect to redis")
)

const (
	ConfigFolder = "config"
)

// BaseConfig is all basic configuration (debug, redis connection and database connection)
type BaseConfig struct {
	Debug    bool
	Author   string
	Redis    RedisCredentials
	Database DatabaseCredentials
}

type RedisCredentials struct {
	Address  string
	Password string
	DB       int
}

type DatabaseCredentials struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
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
	return toml.Unmarshal(c, &cfg)
}

// SetupConfigs with the given configs (+ base config which is available at BaseCfg)
func SetupConfigs(cfgInfo []*ConfigInfo) error {
	err := getBaseConfig(&BaseCfg, DefaultBaseConfig)
	if err != nil {
		return err
	}

	Debug = BaseCfg.Debug
	utils.Author = BaseCfg.Author

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

	_, err = BaseCfg.Redis.Get()
	if err != nil {
		utils.SendAlert("config.go - connection to redis", err.Error())
		return ErrImpossibleToConnectRedis
	}
	return nil
}
