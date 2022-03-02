package botservice

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	entitytype "github.com/oa-dmitriev/telegram-bot/internal/entity-type"
)

func (i *Implementation) Resolve(c *gin.Context) {
	u := tgbotapi.Update{}
	err := c.BindJSON(&u)
	if err != nil {
		log.Println("could not bind json: ", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var chat tgbotapi.Chattable

	ctx := context.Background()
	et := entitytype.New(&u)
	switch {
	case entitytype.IsCommand(et):
		chat, err = i.Command(ctx, u.Message)
	case entitytype.IsMessage(et):
		chat, err = i.Definition(ctx, u.Message, 0)
	case entitytype.IsCallback(et):
		chat, err = i.Callback(ctx, u.CallbackQuery)
	default:
		err = fmt.Errorf("unkown payload")
	}
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := i.Send(ctx, chat); err != nil {
		log.Println("bot send error: ", err)
	}
}

func (i *Implementation) Send(ctx context.Context, chat tgbotapi.Chattable) error {
	if _, err := i.bot.Send(chat); err != nil {
		return err
	}
	return nil
}
