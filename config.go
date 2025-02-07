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
	// GetRedisCredentials returns the RedisCredentials used by the bot
	GetRedisCredentials() *RedisCredentials
	// GetSQLCredentials returns the SQLCredentials used by the bot
	GetSQLCredentials() *SQLCredentials
	// SetDefaultValues set all values of this config to their default ones. THIS IS A DESTRUCTIVE OPERATION!
	SetDefaultValues()
	// Marshal the config, must use toml.Marshal
	Marshal() ([]byte, error)
	// Unmarshal the config, must use toml.Unmarshal
	Unmarshal([]byte) error
}

// SimpleConfig is all basic configuration (debug, redis connection and database connection)
type SimpleConfig struct {
	Debug    bool              `toml:"debug"`
	Author   string            `toml:"author"`
	Redis    *RedisCredentials `toml:"redis"`
	Database *SQLCredentials   `toml:"database"`
}

func (c *SimpleConfig) IsDebug() bool {
	return c.Debug
}

func (c *SimpleConfig) GetAuthor() string {
	return c.Author
}

func (c *SimpleConfig) GetRedisCredentials() *RedisCredentials {
	return c.Redis
}

func (c *SimpleConfig) GetSQLCredentials() *SQLCredentials {
	return c.Database
}

func (c *SimpleConfig) SetDefaultValues() {
	c.Debug = false
	c.Author = "anhgelus"
	c.Redis = &RedisCredentials{}
	c.Redis.SetDefaultValues()
	c.Database = &SQLCredentials{}
	c.Database.SetDefaultValues()
}

func (c *SimpleConfig) Marshal() ([]byte, error) {
	return toml.Marshal(c)
}

func (c *SimpleConfig) Unmarshal(b []byte) error {
	return toml.Unmarshal(b, c)
}

type RedisCredentials struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// SetDefaultValues set all values of these credentials to their default ones. THIS IS A DESTRUCTIVE OPERATION!
func (rc *RedisCredentials) SetDefaultValues() {
	rc.Address = "localhost:6379"
	rc.Password = "password"
	rc.DB = 0
}

type SQLCredentials struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
	Port     int    `toml:"port"`
}

// SetDefaultValues set all values of these credentials to their default ones. THIS IS A DESTRUCTIVE OPERATION!
func (sc *SQLCredentials) SetDefaultValues() {
	sc.Host = "localhost"
	sc.User = "root"
	sc.Password = "root"
	sc.DBName = "bot"
	sc.Port = 5432
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
		utils.SendAlert("config.go - Create file", "File not found, creating a new one.")
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
	return toml.Unmarshal(c, cfg)
}

// SetupConfigs with the given configs (+ base config which is available at BaseCfg)
//
// customBaseConfig is the new type of BaseCfg (if you want to use SimpleConfig, you should pass nil)
func SetupConfigs(customBaseConfig BaseConfig, cfgInfo []*ConfigInfo) error {
	var err error
	if customBaseConfig != nil {
		BaseCfg = customBaseConfig
	} else {
		BaseCfg = &SimpleConfig{}
	}
	err = setupBaseConfig()
	if err != nil {
		return err
	}

	Debug = BaseCfg.IsDebug()
	utils.Author = BaseCfg.GetAuthor()
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
	c, err := BaseCfg.GetRedisCredentials().Connect()
	if err != nil {
		utils.SendAlert("config.go - connection to redis", err.Error())
		return ErrImpossibleToConnectRedis
	}
	_ = c.Close()
	return nil
}
