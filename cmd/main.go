package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"stavki/external/hash"
	"stavki/internal/cache"
	"stavki/internal/database"
	"stavki/internal/log"
	"stavki/internal/rest"
	"stavki/internal/rest/v1"
	"stavki/internal/scheduler"
	"stavki/internal/service"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	DB    database.Config `envPrefix:"DB_"`
	Rest  rest.Config     `envPrefix:"REST_"`
	Redis struct {
		Host     string `env:"HOST"`
		Port     string `env:"PORT"`
		Password string `env:"PASSWORD"`
	} `envPrefix:"REDIS_"`

	HashSalt  string `env:"HASH_SALT"`
	JWTSecret string `env:"JWT_SECRET"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Load only if not production
	if strings.ToLower(os.Getenv("ENV")) != "production" {
		if err := godotenv.Load(".env"); err != nil {
			log.Warn("Error loading .env file: ", err)
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing environment variables: ", err)
	}

	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	txProvider := database.NewTransactionProvider(db)
	userDB := database.NewUserRepository(db)
	jwtDB := database.NewJWTRepository(db)
	eventDB := database.NewEventRepository(db)
	messageDB := database.NewMessageRepository(db)
	betDB := database.NewBetRepository(db)
	marketDB := database.NewMarketRepository(db)
	outcomeDB := database.NewOutcomeRepository(db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Error initializing redis: ", err)
	}

	r := cache.NewCache(rdb)

	authService := service.NewAuth(txProvider, jwtDB, []byte(cfg.JWTSecret))
	userService := service.NewUser(txProvider, userDB, r, hash.NewHasher(cfg.HashSalt), authService)
	eventService := service.NewEvent(txProvider, eventDB, r, *userService)
	chatService := service.NewChatService(messageDB, userService)
	betService := service.NewBetService(betDB)
	marketService := service.NewMarket(txProvider, marketDB, r)
	outcomeService := service.NewOutcome(txProvider, outcomeDB, r)

	jobLogger := log.New()

	s := scheduler.New()
	s.Add(jobLogger, "refresh token clean", 10*time.Second, func(l log.FieldLogger) error {
		l.Info("Cleaning up expired tokens")
		return authService.CleanExpiredTokens(context.Background(), l)
	})

	if err := s.StartAll(); err != nil {
		log.Fatal("Error starting scheduler`s jobs: ", err)
	}

	srv, err := rest.New(&cfg.Rest, log.New()).Init(v1.Config{
		UserService:    userService,
		AuthService:    authService,
		ChatService:    chatService,
		EventService:   eventService,
		BetService:     betService,
		MarketService:  marketService,
		OutcomeService: outcomeService,
	})
	if err != nil {
		log.Fatal("Error initializing server: ", err)
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("listen: ", err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error stopping server")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.StopAll()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := rdb.Close(); err != nil {
			log.WithFields(log.Fields{
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
	log.Info("App exiting")
}
