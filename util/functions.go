package util

import (
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//RandToken returns random string
func RandToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//Get extract data from DefaultQuery -> DefaultPostForm -> GetHeader
func Get(c *gin.Context, key string) string {
	return strings.TrimSpace(c.DefaultQuery(key, c.DefaultPostForm(key, c.GetHeader(key))))
}
