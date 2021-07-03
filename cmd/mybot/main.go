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
	b, err := bot.NewBot(PericBotToken)
	if err != nil {
		log.Fatal(err)
	}
	botHandler, err := handler.NewBotHandler(b, WebHookURL)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.POST("/", botHandler.GetMessage)
	r.Run()
}
