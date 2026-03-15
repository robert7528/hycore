package database

import (
	"github.com/robert7528/hycore/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect opens a GORM connection to the admin database using config.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
}
