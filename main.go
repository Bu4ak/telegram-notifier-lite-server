package main

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func main() {
	r := gin.Default()
	r.Any("/api/send", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"token": getToken(c),
		})
	})
	r.Run()
}

type Req struct {
	Token string `json:"token"`
}

func getToken(c *gin.Context) string {
	key := "token"
	var req Req
	token := c.DefaultQuery(key, c.DefaultPostForm(key, c.GetHeader(key)))
	if token == "" {
		if err := c.BindJSON(&req); err == nil {
			token = req.Token
		}
	}
	return strings.TrimSpace(token)
}
