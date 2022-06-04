package botservice

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oa-dmitriev/telegram-bot/internal/domain"
	"github.com/oa-dmitriev/telegram-bot/internal/markup"
)

func (i *Implementation) Command(ctx context.Context, msg *tgbotapi.Message) (tgbotapi.Chattable, error) {
	if msg.Command() == "vocab" {
		data, err := i.vocabRepo.GetList(ctx, int64(msg.From.ID), pageLen, 0)
		if err != nil {
			return nil, err
		}

		sendMsg := tgbotapi.NewMessage(msg.Chat.ID, "")
		sendMsg.ReplyToMessageID = msg.MessageID
		sendMsg.ParseMode = "markdown"

		if len(data) == 0 {
			sendMsg.Text = "Your vocabulary is empty"
			return sendMsg, nil
		}

		domainData := domain.ConvertDBVocabToDomainData(data)

		sendMsg.Text = domain.ToString(domainData)

		mark := markup.New()
		if len(domainData) == pageLen {
			mark = mark.WithNext(0)
			sendMsg.ReplyMarkup = mark.InlineKeyboardMarkup
		}

		return sendMsg, nil
	}

	if msg.Command() == "delete" {
		userID := msg.From.ID
		if err := i.vocabRepo.Delete(ctx, int64(userID)); err != nil {
			return nil, err
		}
		sendMsg := tgbotapi.NewMessage(msg.Chat.ID, "Your vocabulary has been deleted")
		sendMsg.ReplyToMessageID = msg.MessageID
		sendMsg.ParseMode = "markdown"
		return sendMsg, nil
	}
	return nil, fmt.Errorf("uknown command")
}
