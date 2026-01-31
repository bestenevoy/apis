package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"wrzapi/internal/httpclient"
	"wrzapi/internal/pageinfo"
)

func PageInfo(c *gin.Context) {
	raw := c.Query("url")
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing url"})
		return
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url must be http or https"})
		return
	}

	client := httpclient.New()
	body, finalURL, err := client.FetchHTML(c.Request.Context(), parsed.String())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	meta := pageinfo.ParseHTML(body, finalURL)
	c.JSON(http.StatusOK, meta)
}
