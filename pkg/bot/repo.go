package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotRepo struct {
	Bot *tgbotapi.BotAPI
}

func NewBotRepo(token string, webHookUrl string) (*BotRepo, error) {
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	_, err = b.SetWebhook(tgbotapi.NewWebhook(webHookUrl))
	if err != nil {
		return nil, err
	}
	return &BotRepo{Bot: b}, nil
}

func (repo *BotRepo) Edit(msg *tgbotapi.Message) {
	// msg := tgbotapi.NewMessage(u.Message.Chat.ID, u.Message.Text)
	// msg.ReplyToMessageID = u.Message.MessageID
	// bh.Bot.Send(msg)
	// log.Println("GETALL WORKING with bindJSON")
}

func (repo *BotRepo) Message(msg *tgbotapi.Message) error {
	tokens := strings.Fields(msg.Text)
	if len(tokens) < 1 {
		return nil
	}
	if cmd, ok := cmds[tokens[0]]; ok {
		return cmd(repo, msg, tokens[1:])
	}
	return fmt.Errorf("no valid command")
}
