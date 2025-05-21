package gokord

import (
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
