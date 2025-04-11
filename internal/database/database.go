package database

import (
	"fmt"

	"stavki/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Config struct {
		// Host is the database host.
		Host string `env:"HOST"`
		// Port is the database port.
		Port string `env:"PORT"`
		// User is the database user.
		User string `env:"USERNAME"`
		// Password is the database password.
		Password string `env:"PASSWORD"`
		// Name is the database name.
		Name string `env:"DATABASE"`
		// Schema is the database schema.
		Schema string `env:"SCHEMA"`
	}
)

func Connect(cfg Config) (*gorm.DB, error) {
	// Create the database connection.
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s search_path=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Password, cfg.Schema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Migrate the database.
	if err := db.AutoMigrate(model.Migrate()...); err != nil {
		return nil, err
	}

	return db, nil
}
