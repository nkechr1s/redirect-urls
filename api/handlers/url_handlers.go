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
