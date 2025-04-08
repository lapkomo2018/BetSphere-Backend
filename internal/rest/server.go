package rest

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"stavki/internal/rest/v1"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	r.Use(gin.LoggerWithWriter(logrus.StandardLogger().Writer()))

	return &Server{
		gin: r,
		cfg: cfg,
	}
}

func (s *Server) Init(cfgV1 v1.Config) *http.Server {
	v1.New(cfgV1).
		Init(s.gin.Group("/v1"))

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Port),
		Handler: s.gin.Handler(),
	}
}
