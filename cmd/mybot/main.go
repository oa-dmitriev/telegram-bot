package main

import (
	"log"
	"mybot/pkg/handler"
	"mybot/pkg/mybot"
	"time"
)

var (
	PericBotToken = "1887809470:AAGvUYu2u-H1DQd16Si0Dpag3-sME_2HjwY"
	WebHookURL    = "https://peric-telegram-bot.herokuapp.com"
)

func main() {
	bot, err := mybot.NewBot(PericBotToken)
	if err != nil {
		log.Fatal(err)
	}
	_, err = handler.NewBotHandler(bot, WebHookURL)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(20 * time.Minute)
}
