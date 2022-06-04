package botservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	entitytype "github.com/oa-dmitriev/telegram-bot/internal/entity-type"
	"github.com/oa-dmitriev/telegram-bot/internal/repository"
	"github.com/oa-dmitriev/telegram-bot/internal/repository/user"
)

var congratulationsMsg = `
	Sup dude! This is your first time using this awesome bot!
	I hope you'll have a great time!
`

func (i *Implementation) Resolve(c *gin.Context) {
	u := tgbotapi.Update{}
	err := c.BindJSON(&u)
	if err != nil {
		log.Println("could not bind json: ", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx := context.Background()
	et := entitytype.New(&u)

	var chat tgbotapi.Chattable
	var firstMsg *tgbotapi.MessageConfig

	if entitytype.IsMessage(et) {
		firstMsg, err = i.CreateUserIfNotExists(ctx, u.Message)
		if err != nil {
			log.Printf("error creating user: %s\n", err.Error())
			return
		}
	}

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
	}

	if err := i.Send(ctx, chat, firstMsg); err != nil {
		log.Println("bot send error: ", err)
	}
}

func (i *Implementation) Send(ctx context.Context, chat tgbotapi.Chattable, firstMsg *tgbotapi.MessageConfig) error {
	if chat == nil {
		return fmt.Errorf("empty response")
	}

	if firstMsg != nil {
		i.bot.Send(firstMsg)
	}

	_, err := i.bot.Send(chat)

	return err
}

func (i *Implementation) CreateUserIfNotExists(ctx context.Context, msg *tgbotapi.Message) (*tgbotapi.MessageConfig, error) {
	if msg == nil {
		return nil, errors.New("message is empty")
	}
	if msg.From == nil {
		return nil, errors.New("user is empty")
	}

	var err error
	if _, err = i.userRepo.GetUser(ctx, int64(msg.From.ID)); errors.Is(err, user.ErrUserNotFound) {
		u := &repository.DBUser{
			ID:        int64(msg.From.ID),
			Username:  msg.From.UserName,
			FirstName: msg.From.FirstName,
			LastName:  msg.From.LastName,
			ChatID:    msg.Chat.ID,
		}
		congrats := tgbotapi.NewMessage(msg.Chat.ID, congratulationsMsg)
		return &congrats, i.userRepo.Add(ctx, u)
	}
	return nil, err
}
