package main

import (
	"net/http"
	"redirectUrls/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	err := handlers.InitializeDB()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.GET("/urls", handlers.GetUrls)
	router.GET("/urls/:id", handlers.GetUrlByID)
	router.POST("/urls", handlers.CreateUrl)
	router.DELETE("/urls/:id", handlers.DeleteUrlByID)
	router.PATCH("/urls/:id", handlers.PatchUrlByID)

	router.POST("/generate-nginx-config", func(c *gin.Context) {
		handlers.FetchUrls()
		c.JSON(http.StatusOK, gin.H{"message": "NGINX configuration file generated successfully"})
	})
	router.Run("localhost:8080")
}
