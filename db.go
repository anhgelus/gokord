package gokord

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB used
var DB *gorm.DB

// DataBase is an interface with basic methods to load and save data
type DataBase interface {
	Load() error // Load data from the database
	Save() error // Save data into the database
}

type BotData struct {
	gorm.Model
	Version string `gorm:"version"`
	Name    string `gorm:"name"`
}

func (b *BotData) Load() error {
	return DB.FirstOrCreate(b).Error
}

func (b *BotData) Save() error {
	return DB.Save(b).Error
}

// Connect to the postgres database using the given config.DatabaseCredentials
func (dc *DatabaseCredentials) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dc.generateDsn()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// generateDsn for the connection to postgres using the given config.DatabaseCredentials
func (dc *DatabaseCredentials) generateDsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Paris",
		dc.Host, dc.User, dc.Password, dc.DBName, dc.Port,
	)
}
