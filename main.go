package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	gin.SetMode(os.Getenv("GIN_MODE"))
}

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
	go listenUpdates()
	r.Run()
}

func get(c *gin.Context, key string) string {
	return strings.TrimSpace(c.DefaultQuery(key, c.DefaultPostForm(key, c.GetHeader(key))))
}

func listenUpdates() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug, _ = strconv.ParseBool(os.Getenv("TELEGRAM_BOT_DEBUG"))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		go bot.Send(msg)
	}
}
