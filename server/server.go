package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/arief-hidayat/go-video-api/controllers"
	"github.com/arief-hidayat/go-video-api/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/caarlos0/env/v6"
)

type VideoAPI struct {
	config config
	e      *echo.Echo
}

type AdminAPI struct {
	config config
}

func (a *AdminAPI) GetAppVersion(c echo.Context) (err error) {
	return c.String(http.StatusOK, a.config.AppVersion)
}

// VideoAPI Instance of Echo
func NewServer() *VideoAPI {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &VideoAPI{
		e:      echo.New(),
		config: cfg,
	}
}

type config struct {
	AppVersion          string          `env:"APP_VERSION" envDefault:"latest"`
	DbHost              string          `env:"DB_HOST" envDefault:"127.0.0.1"`
	DbPort              int             `env:"DB_PORT" envDefault:"5432"`
	DbUser              string          `env:"DB_USER" envDefault:"vod123"`
	DbPwd               string          `env:"DB_PWD"`
	DbName              string          `env:"DB_NAME" envDefault:"vod"`
	DbNoSSL             bool            `env:"DB_NO_SSL" envDefault:"false"`
	MaxOpenConnections  int             `env:"MAX_OPEN_CONNECTIONS" envDefault:"5"`
	ConnTimeout         time.Duration   `env:"CONNECTION_TIMEOUT" envDefault:"20s"`
}

// Start server functionality
func (s *VideoAPI) Start(port string) {
	// logger
	s.e.Use(middleware.Logger())
	// recover
	s.e.Use(middleware.Recover())
	//CORS
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	s.e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
        Timeout: s.config.ConnTimeout,
    }))
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		s.config.DbUser, url.QueryEscape(s.config.DbPwd), s.config.DbHost, s.config.DbPort, s.config.DbName)
	if s.config.DbNoSSL {
		connStr = fmt.Sprintf("%s?sslmode=disable", connStr)
	}
	err := models.InitDB(connStr, s.config.MaxOpenConnections)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	adminApi := &AdminAPI{config: s.config}
	// videos endpoint
	s.e.GET("/", adminApi.GetAppVersion)
	s.e.File("/favicon.ico", "assets/favicon.png")
	s.e.GET("/videos", controllers.GetVideos)
	// Start Server
// 	s.e.Logger.Fatal(s.e.Start(port))

    // https://echo.labstack.com/cookbook/graceful-shutdown/
	// Start server
	go func() {
		if err := s.e.Start(port); err != nil && err != http.ErrServerClosed {
			s.e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	models.CloseDB()
	if err := s.e.Shutdown(ctx); err != nil {
		s.e.Logger.Fatal(err)
	}
}

// Close server functionality
func (s *VideoAPI) Close() {
    models.CloseDB()
	s.e.Close()
}
