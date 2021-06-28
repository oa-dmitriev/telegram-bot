package mybot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func NewBot(token string) (*tgbotapi.BotAPI, error) {
	return tgbotapi.NewBotAPI(token)
}
