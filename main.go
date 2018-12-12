package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func main() {
	r := gin.New()
	if gin.IsDebugging() {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	r.Any("/api/send", func(c *gin.Context) {
		var req struct {
			Token   string `json:"token"`
			Message string `json:"message"`
		}
		req.Token = get(c, "token")
		req.Message = get(c, "message")

		if c.ContentType() == gin.MIMEJSON {
			c.BindJSON(&req)
		}
		c.JSON(http.StatusOK, req)
	})
	r.Run()
}

func get(c *gin.Context, key string) string {
	return strings.TrimSpace(c.DefaultQuery(key, c.DefaultPostForm(key, c.GetHeader(key))))
}
