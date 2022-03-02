package botservice

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oa-dmitriev/telegram-bot/internal/markup"
	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

func (i *Implementation) Callback(ctx context.Context, cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	chatID, msgID, userID := cb.Message.Chat.ID, cb.Message.MessageID, cb.From.ID
	word := cb.Message.ReplyToMessage.Text

	sendMsg := tgbotapi.NewEditMessageText(chatID, msgID, "")
	sendMsg.ParseMode = "markdown"

	if cb.Data == "add" {
		if err := i.save(ctx, int64(userID), word); err != nil {
			sendMsg.Text = err.Error()
			return nil, fmt.Errorf("failed to save the word [%s] for the user [%d]", word, userID)
		}
		sendMsg.Text = fmt.Sprintf("*%s* added\n/vocab - show your vocabulary", word)
		return sendMsg, nil
	}

	pageNum, err := strconv.Atoi(cb.Data)
	if err != nil {
		return nil, err
	}

	if cb.Message.ReplyToMessage.Text == "/vocab" {
		offset := int64(pageNum) * pageLen
		DBVocab, err := i.getListFromRepo(ctx, offset, int64(userID))
		if err != nil {
			return nil, err
		}
		mark := markup.New().WithPrev(pageNum)
		if len(DBVocab) == pageLen {
			mark = mark.WithNext(pageNum)
		}
		sendMsg.ReplyMarkup = mark.InlineKeyboardMarkup
		return sendMsg, nil
	}

	msg, err := i.Definition(ctx, cb.Message, pageNum)
	if err != nil {
		return nil, err
	}
	if edditMsg, ok := msg.(tgbotapi.MessageConfig); ok {
		if inlineMarkup, success := (edditMsg.ReplyMarkup).(*tgbotapi.InlineKeyboardMarkup); success {
			sendMsg.ReplyMarkup = inlineMarkup
			return sendMsg, nil
		}
		log.Println("failed casting to InlineKeyboardMarkup")
		return sendMsg, nil
	}

	return nil, fmt.Errorf("failed casting to MessageConfig")
}

func (i *Implementation) getListFromRepo(
	ctx context.Context,
	offset, userID int64,
) ([]*repository.DBVocabulary, error) {
	DBVocab, err := i.vocabRepo.GetList(ctx, userID, pageLen, offset)
	if err != nil {
		return nil, err
	}
	return DBVocab, nil
}

func (i *Implementation) save(ctx context.Context, userID int64, word string) error {
	defs, err := i.FetchData(ctx, word)
	if err != nil {
		return err
	}

	if len(defs) == 0 {
		return fmt.Errorf("[%s] defintion not found", word)
	}

	return i.vocabRepo.Add(ctx, &repository.DBVocabulary{
		UserID:     userID,
		Word:       word,
		Definition: defs[0],
	})
}
