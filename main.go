package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var tokens = make(map[int64]string)
var channels = make(map[string]int64)
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var bot *tgbotapi.BotAPI
var err error

func init() {
	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	gin.SetMode(os.Getenv("GIN_MODE"))

	if bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN")); err != nil {
		log.Panic(err)
	}
	bot.Debug, _ = strconv.ParseBool(os.Getenv("TELEGRAM_BOT_DEBUG"))
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
		chatID, exists := channels[req.Token]
		if !exists {
			c.Status(http.StatusUnauthorized)
			return
		}
		go bot.Send(tgbotapi.NewMessage(chatID, req.Message))

		c.Status(http.StatusOK)
	})
	go listenUpdates()
	r.Run()
}

func get(c *gin.Context, key string) string {
	return strings.TrimSpace(c.DefaultQuery(key, c.DefaultPostForm(key, c.GetHeader(key))))
}

func listenUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		token, exists := tokens[update.Message.Chat.ID]
		if !exists {
			token = randToken(30)
			tokens[update.Message.Chat.ID] = token
			channels[token] = update.Message.Chat.ID
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This channel token: `"+token+"`")
		msg.ParseMode = "markdown"
		log.Println(tokens)
		log.Println(channels)
		go bot.Send(msg)
	}
}

func randToken(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
