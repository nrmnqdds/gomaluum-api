package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nrmnqdds/gomaluum-api/services/auth"
	"github.com/nrmnqdds/gomaluum-api/services/scraper"
	"net/http"
)

func main() {
	r := gin.Default()

	// Apply middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/api/v1/login", func(c *gin.Context) {

		type LoginSchema struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		var json LoginSchema
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		auth.LoginUser(c, json.Username, json.Password)

	})

	r.GET("/api/v1/schedule", func(c *gin.Context) {
		cookie, err := c.Cookie("MOD_AUTH_CAS")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if cookie == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please login first"})
			return
		}

		scraper.ScheduleScraperService(cookie)

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	fmt.Println("Server is running on port 8080")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
