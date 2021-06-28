package main

import (
	"io/ioutil"
	"log"
	"mybot/pkg/handler"
	"mybot/pkg/mybot"
	"net/http"
	"os"
)

var (
	PericBotToken = "1887809470:AAGvUYu2u-H1DQd16Si0Dpag3-sME_2HjwY"
	WebHookURL    = "https://peric-telegram-bot.herokuapp.com"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	log.Println("In: GETALL")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("ERR: ", err)
		return
	}
	log.Println(string(body))
}

func main() {
	bot, err := mybot.NewBot(PericBotToken)
	if err != nil {
		log.Fatal(err)
	}
	_, err = handler.NewBotHandler(bot, WebHookURL)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", GetAll)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
