package main

import (
	"redirectUrls/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/urls", handlers.GetUrls)
	router.GET("/urls/:id", handlers.GetUrlByID)
	router.Run("localhost:8080")
}
