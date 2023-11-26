package handlers

import (
	"database/sql"
	"net/http"
	"redirectUrls/api/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetUrls(c *gin.Context) {
	rows, err := db.Query("SELECT id, currentUrl, redirectUrl FROM urls")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.CurrentUrl, &url.RedirectUrl)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		urls = append(urls, url)
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": urls})
}

func GetUrlByID(c *gin.Context) {
	id := c.Param("id")

	var url models.URL
	err := db.QueryRow("SELECT id, currentUrl, redirectUrl FROM urls WHERE id = ?", id).
		Scan(&url.ID, &url.CurrentUrl, &url.RedirectUrl)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "URL not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": url})
}
