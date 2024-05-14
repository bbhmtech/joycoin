package joycoin

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDatabase(cfg *Config) error {
	db, err := gorm.Open(sqlite.Open(cfg.DatabseConnectionString))
}
