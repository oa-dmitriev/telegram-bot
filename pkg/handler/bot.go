package handler

import (
	"log"
	"mybot/pkg/bot"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotHandler struct {
	BotRepo *bot.BotRepo
}

func (h *BotHandler) GetMessage(c *gin.Context) {
	u := tgbotapi.Update{}
	err := c.BindJSON(&u)
	log.Printf("%#v\n", u)
	if err != nil {
		log.Println("ERR: ", err)
		return
	}

	if u.Message != nil {
		h.BotRepo.Message(u.Message)
		return
	}
}
