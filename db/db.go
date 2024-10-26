package db

import (
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	slog.Info("Initializing database")
	database, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db = database
	migrate()
	slog.Info("Database initialized")
}

func migrate() {
	slog.Info("Migrating database")
	db.AutoMigrate(&RateObject{})
	db.AutoMigrate(&Provider{})
	slog.Info("Database migration complete")
}
