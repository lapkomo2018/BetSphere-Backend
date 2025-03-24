package database

import (
	"context"
	"fmt"

	"stavki/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Database interface {
		Transaction() *Transaction
		TransactionWithContext(ctx context.Context) *Transaction
	}

	database struct {
		db *gorm.DB
	}

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

func New(cfg Config) (Database, error) {
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

	return &database{
		db: db,
	}, nil
}

func (db *database) Transaction() *Transaction {
	return &Transaction{
		DB: db.db.Begin(),
	}
}

func (db *database) TransactionWithContext(ctx context.Context) *Transaction {
	return &Transaction{
		DB: db.db.Begin().WithContext(ctx),
	}
}
