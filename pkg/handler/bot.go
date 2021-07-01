package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotHandler struct {
	Bot *tgbotapi.BotAPI
}

func NewBotHandler(bot *tgbotapi.BotAPI, WebHookURL string) (*BotHandler, error) {
	bh := BotHandler{bot}
	_, err := bh.Bot.SetWebhook(tgbotapi.NewWebhook(WebHookURL))
	if err != nil {
		return nil, err
	}
	return &bh, nil
}

func (b *BotHandler) GetAll(c *gin.Context) {
	log.Println("GETALL WORKING")
	body, err := c.GetRawData()
	if err != nil {
		log.Println("ERR: ", err)
		return
	}
	log.Println(string(body))
}
