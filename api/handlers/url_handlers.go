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
	c.IndentedJSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": urls})
}

func GetUrlByID(c *gin.Context) {
	id := c.Param("id")
	for _, url := range urls {
		if url.ID == id {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": url})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "URL not found"})
}
