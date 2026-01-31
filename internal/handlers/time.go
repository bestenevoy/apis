package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func Time(c *gin.Context) {
	now := time.Now().UTC()
	resp := gin.H{
		"utc":   now.Format(time.RFC3339),
		"unix":  now.Unix(),
		"local": time.Now().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, resp)
}
