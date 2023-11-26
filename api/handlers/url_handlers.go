package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"redirectUrls/api/models"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetUrls(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM urls")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.CurrentUrl, &url.RedirectUrl, &url.CreatedAt, &url.UpdatedAt)
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
	err := db.QueryRow("SELECT * FROM urls WHERE id = ?", id).
		Scan(&url.ID, &url.CurrentUrl, &url.RedirectUrl, &url.CreatedAt, &url.UpdatedAt)

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

	var existingID string
	err := db.QueryRow("SELECT id FROM urls WHERE currentUrl = ?", newURL.CurrentUrl).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"status": http.StatusConflict, "message": "URL already exists with the given currentUrl"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO urls (currentUrl, redirectUrl, createdAt, updatedAt) VALUES (?, ?, ?, ?)", newURL.CurrentUrl, newURL.RedirectUrl, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	newID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	newURL.ID = fmt.Sprintf("%d", newID)

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Url created successfully"})
}

func DeleteUrlByID(c *gin.Context) {
	id := c.Param("id")

	var existingID string
	err := db.QueryRow("SELECT id FROM urls WHERE id = ?", id).Scan(&existingID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "URL not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	_, err = db.Exec("DELETE FROM urls WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "URL deleted successfully"})
}

func PatchUrlByID(c *gin.Context) {
	id := c.Param("id")

	var existingID string
	err := db.QueryRow("SELECT id FROM urls WHERE id = ?", id).Scan(&existingID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "URL not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	var existingURL models.URL
	err = db.QueryRow("SELECT * FROM urls WHERE id = ?", id).
		Scan(&existingURL.ID, &existingURL.CurrentUrl, &existingURL.RedirectUrl, &existingURL.CreatedAt, &existingURL.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&existingURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE urls SET redirectUrl = ?, updatedAt = ? WHERE id = ?", existingURL.RedirectUrl, time.Now().Unix(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	err = db.QueryRow("SELECT * FROM urls WHERE id = ?", id).
		Scan(&existingURL.ID, &existingURL.CurrentUrl, &existingURL.RedirectUrl, &existingURL.CreatedAt, &existingURL.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Url updated successfully"})
}
