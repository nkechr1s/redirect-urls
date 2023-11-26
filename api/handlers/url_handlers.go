package handlers

import (
	"net/http"
	"redirectUrls/api/models"

	"github.com/gin-gonic/gin"
)

// TODO REPLACE WITH REAL DATA
var urls = []models.URL{
	{ID: "1", CurrentUrl: "/blackfriday", RedirectUrl: "/cybermonday"},
	{ID: "2", CurrentUrl: "/blackfriday", RedirectUrl: "/cybermonday"},
}

func GetUrls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, urls)
}

func GetUrlByID(c *gin.Context) {
	id := c.Param("id")
	for _, a := range urls {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "url not found"})
}
