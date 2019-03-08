package main

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"telegram-notifier-server/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var bot *tgbotapi.BotAPI
var db *sql.DB
var router *gin.Engine

func init() {
	var err error
	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	gin.SetMode(os.Getenv("GIN_MODE"))
	if gin.Mode() == "release" {
		logFile, _ := os.Create("error.log")
		gin.DefaultErrorWriter = logFile
	}
	router = gin.New()
	router.Use(gin.Recovery())
	if gin.IsDebugging() {
		router.Use(gin.Logger())
	}

	db, err = sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN")); err != nil {
		log.Fatal(err)
	}
	bot.Debug, _ = strconv.ParseBool(os.Getenv("TELEGRAM_BOT_DEBUG"))
}

func main() {
	router.Any("/api/send", func(c *gin.Context) {
		var req struct {
			Token   string `json:"token"`
			Message string `json:"message"`
		}
		req.Token = util.Get(c, "token")
		req.Message = util.Get(c, "message")

		if c.ContentType() == gin.MIMEJSON {
			c.BindJSON(&req)
		}
		chatID := getChatIdByToken(req.Token)

		if chatID == 0 {
			c.Status(http.StatusUnauthorized)
			return
		}
		go bot.Send(tgbotapi.NewMessage(chatID, req.Message))

		c.Status(http.StatusOK)
	})
	go listenUpdates()
	router.Run()
}

func listenUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		token := getTokenById(update.Message.Chat.ID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This channel token: `"+token+"`")
		msg.ParseMode = "markdown"
		go bot.Send(msg)
	}
}

func getChatIdByToken(token string) (chatID int64) {
	q := "select chat_id from users where token = $1"
	if e := db.QueryRow(q, token).Scan(&chatID); e != nil && e.Error() != "sql: no rows in result set" {
		log.Panic(e.Error())
	}
	return
}

func getTokenById(chatID int64) (token string) {
	if e := db.QueryRow("select token from users where chat_id = $1", chatID).Scan(&token); e != nil {
		if e.Error() == "sql: no rows in result set" {
			token = util.RandToken(32)
			q := "insert into users (chat_id, token, created_at) values ($1, $2, now())"
			if _, err := db.Exec(q, chatID, token); err != nil {
				log.Panic(e.Error())
			}
		} else {
			log.Panic(e.Error())
		}
	}
	return
}
