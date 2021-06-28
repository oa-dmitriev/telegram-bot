package main

import (
	"log"
	"mybot/pkg/handler"
	"mybot/pkg/mybot"
	"os"
	"time"
)

var (
	PericBotToken = "1887809470:AAGvUYu2u-H1DQd16Si0Dpag3-sME_2HjwY"
	WebHookURL    = "0.0.0.0"
)

func main() {
	bot, err := mybot.NewBot(PericBotToken)
	if err != nil {
		log.Fatal(err)
	}
	_, err = handler.NewBotHandler(bot, WebHookURL+":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(20 * time.Minute)
}
