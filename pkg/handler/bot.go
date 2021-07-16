package handler

import (
	"fmt"
	"log"
	"mybot/pkg/bot"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotHandler struct {
	BotRepo *bot.BotRepo
}

func (h *BotHandler) ProcessMessage(c *gin.Context) {
	u := tgbotapi.Update{}
	err := c.BindJSON(&u)
	log.Printf("%#v\n", u)
	if err != nil {
		log.Println("\n\n\n\nERR: ", err)
		return
	}

	var chat tgbotapi.Chattable

	if u.Message != nil && u.Message.Command() != "" {
		chat, err = h.BotRepo.Command(u.Message)
	} else if u.Message != nil {
		chat, err = h.BotRepo.Definition(u.Message)
	} else if u.CallbackQuery != nil {
		chat, err = h.BotRepo.CallBack(u.CallbackQuery)
	} else {
		err = fmt.Errorf("Uknown payload")
	}
	if err != nil {
		log.Println("\n\n\n\nERROR: ", err)
		return
	}
	h.BotRepo.Bot.Send(chat)
}
