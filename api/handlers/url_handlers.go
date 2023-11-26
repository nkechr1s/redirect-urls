package handlers

import (
	"database/sql"
	"fmt"
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

func CreateUrl(c *gin.Context) {
	var newURL models.URL
	if err := c.ShouldBindJSON(&newURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO urls (currentUrl, redirectUrl) VALUES (?, ?)", newURL.CurrentUrl, newURL.RedirectUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	// Get the auto-incremented ID of the newly inserted URL
	newID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	// Update the ID of the newURL before sending it back in the response
	newURL.ID = fmt.Sprintf("%d", newID)

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "data": newURL})
}
