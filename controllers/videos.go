package controllers

import (
	"log"
	"net/http"

	"github.com/arief-hidayat/go-video-api/models"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
)

// GetVideos ...
func GetVideos(c echo.Context) (err error) {
	queryStr := c.Request().FormValue("q")
	videos, err := models.GetVideos(queryStr)
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, videos)
}
