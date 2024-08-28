package main

import (
	"github.com/labstack/echo/v4"
	"github.com/nrmnqdds/gomaluum-api/handlers/auth"
	_ "github.com/nrmnqdds/gomaluum-api/handlers/scraper"
	"net/http"
)

func main() {
	e := echo.New()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.POST("/api/v1/login", auth.LoginUser)

	// e.GET("/api/v1/schedule", func(c *gin.Context) {
	// 	cookie, err := c.Cookie("MOD_AUTH_CAS")
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	//
	// 	if cookie == "" {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Please login first"})
	// 		return
	// 	}
	//
	// 	scraper.ScheduleScraperService(cookie)
	//
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "success",
	// 	})
	// })

	e.Logger.Fatal(e.Start(":1323"))
}
