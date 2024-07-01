package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nrmnqdds/gomaluum-api/services/auth"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/login", func(c *gin.Context) {

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

	fmt.Println("Server is running on port 8080")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
