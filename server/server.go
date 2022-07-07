package server

import (
	"fmt"
	"go-video-api/controllers"
	"go-video-api/models"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/caarlos0/env/v6"
)

type VideoAPI struct {
	config config
	e      *echo.Echo
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
	DbHost  string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DbPort  int    `env:"DB_PORT" envDefault:"5432"`
	DbUser  string `env:"DB_USER" envDefault:"vod123"`
	DbPwd   string `env:"DB_PWD"`
	DbName  string `env:"DB_NAME" envDefault:"vod"`
	DbNoSSL bool   `env:"DB_NO_SSL" envDefault:"false"`
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

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		s.config.DbUser, s.config.DbPwd, s.config.DbHost, s.config.DbPort, s.config.DbName)
	if s.config.DbNoSSL {
		connStr = fmt.Sprintf("%s?sslmode=disable", connStr)
	}
	err := models.InitDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	// videos endpoint
	s.e.GET("/videos", controllers.GetVideos)
	// Start Server
	s.e.Logger.Fatal(s.e.Start(port))
}

// Close server functionality
func (s *VideoAPI) Close() {
	s.e.Close()
}
