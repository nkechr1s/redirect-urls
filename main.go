package main

import (
	"os"
	"redirectUrls/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	err := handlers.InitializeDB()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	server := os.Getenv("SERVER_NAME")
	router.GET("/urls", handlers.GetUrls)
	router.GET("/urls/:id", handlers.GetUrlByID)
	router.POST("/urls", handlers.CreateUrl)
	router.DELETE("/urls/:id", handlers.DeleteUrlByID)
	router.PATCH("/urls/:id", handlers.PatchUrlByID)

	router.POST("/generate-nginx-config", handlers.GenerateNginxConfig)
	router.Run(server)
}
