package main

import (
	"log"
	"mybot/pkg/bot"
	"mybot/pkg/handler"

	"github.com/gin-gonic/gin"
)

var (
	PericBotToken = "1887809470:AAGvUYu2u-H1DQd16Si0Dpag3-sME_2HjwY"
	WebHookURL    = "https://peric-telegram-bot.herokuapp.com"
)

func main() {
	botRepo, err := bot.NewBotRepo(PericBotToken, WebHookURL)
	if err != nil {
		log.Fatal(err)
	}
	handler := handler.BotHandler{
		BotRepo: botRepo,
	}
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.POST("/", handler.GetMessage)
	r.Run()
}
