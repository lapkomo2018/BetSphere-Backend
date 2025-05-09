package rest

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"stavki/internal/log"
	"stavki/internal/rest/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type (
	Server struct {
		gin *gin.Engine
		cfg *Config
	}

	Config struct {
		Port int `env:"PORT" envDefault:"8080"`
	}
)

func New(cfg *Config) *Server {
	if strings.ToLower(os.Getenv("ENV")) == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())

	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	r.Use(gin.LoggerWithWriter(logger.Writer()))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return &Server{
		gin: r,
		cfg: cfg,
	}
}

func (s *Server) Init(cfgV1 v1.Config) (*http.Server, error) {
	v1Handler, err := v1.New(cfgV1)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to create v1 handler"))
	}
	v1Handler.Init(s.gin.Group("/v1"))

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Port),
		Handler: s.gin.Handler(),
	}, nil
}
