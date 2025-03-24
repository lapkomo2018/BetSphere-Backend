package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"stavki/external/hash"
	"stavki/internal/database"
	"stavki/internal/rest"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DB    database.Config `envPrefix:"DB_"`
	Rest  rest.Config     `envPrefix:"REST_"`
	Redis struct {
		Host     string `env:"HOST"`
		Port     string `env:"PORT"`
		Password string `env:"PASSWORD"`
	} `envPrefix:"REDIS_"`

	HashSalt string `env:"HASH_SALT"`
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Load only if not production
	if strings.ToLower(os.Getenv("ENV")) != "production" {
		if err := godotenv.Load(".env"); err != nil {
			logrus.Warn("Error loading .env file: ", err)
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatal("Error parsing environment variables: ", err)
	}

	_, err := database.New(cfg.DB)
	if err != nil {
		logrus.Fatal("Error initializing database: ", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logrus.Fatal("Error initializing redis: ", err)
	}

	_ = hash.NewHasher(cfg.HashSalt)

	srv := rest.New(&cfg.Rest).Init()

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatal("listen: ", err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error stopping server")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := rdb.Close(); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("Error stopping redis")
		}
	}()

	// Wait for wg done
	go func() {
		wg.Wait()
		cancel()
	}()

	ctx.Done()
	log.Println("App exiting")
}
