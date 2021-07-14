package bot

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotRepo struct {
	Bot    *tgbotapi.BotAPI
	RedCon *redis.Client
}

func InitRedis() (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}
	conn := redis.NewClient(opt)
	status := conn.Ping()
	if status.Err() != nil {
		return nil, status.Err()
	}
	return conn, nil
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
	conn, _ := InitRedis()
	return &BotRepo{Bot: b, RedCon: conn}, nil
}

func (repo *BotRepo) Message(msg *tgbotapi.Message) error {
	tokens := strings.Fields(msg.Text)
	if len(tokens) < 2 {
		return nil
	}
	cmd, ok := cmds[tokens[0]]
	if !ok {
		return fmt.Errorf("no valid command")
	}
	data, err := cmd(repo, tokens[1:])
	if err != nil {
		return fmt.Errorf("server error")
	}
	if len(data) == 0 {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, "No definition found")
		newMsg.ReplyToMessageID = msg.MessageID
		repo.Bot.Send(newMsg)
		return nil
	}
	dataToSend := GetPage(data, 0)
	newMsg := tgbotapi.NewMessage(
		msg.Chat.ID,
		"•  "+strings.Join(dataToSend, "\n•  "),
	)
	if isDataLeft(data, 0) {
		newMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Next", "1"),
			),
		)
	}
	newMsg.ReplyToMessageID = msg.MessageID
	repo.Bot.Send(newMsg)
	return fmt.Errorf("no valid command")
}

func (repo *BotRepo) CallBackQuery(cb *tgbotapi.CallbackQuery) error {
	msg := cb.Message
	num, err := strconv.Atoi(cb.Data)
	if err != nil {
		return fmt.Errorf("uknown callbackquery data")
	}
	tokens := strings.Fields(cb.Message.ReplyToMessage.Text)
	if len(tokens) < 2 {
		return nil
	}
	cmd, ok := cmds[tokens[0]]
	if !ok {
		return fmt.Errorf("no valid command")
	}
	data, err := cmd(repo, tokens[1:])
	if err != nil {
		return fmt.Errorf("server error")
	}

	dataToSend := GetPage(data, num)
	newMsg := tgbotapi.NewEditMessageText(
		msg.Chat.ID,
		msg.MessageID,
		"•  "+strings.Join(dataToSend, "\n•  "),
	)
	if num == 0 {
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Next", strconv.Itoa(num+1)),
			),
		)
		newMsg.ReplyMarkup = &markup
	} else if isDataLeft(data, num) {
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Prev", strconv.Itoa(num-1)),
				tgbotapi.NewInlineKeyboardButtonData("Next", strconv.Itoa(num+1)),
			),
		)
		newMsg.ReplyMarkup = &markup
	} else {
		markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Prev", strconv.Itoa(num-1)),
			),
		)
		newMsg.ReplyMarkup = &markup
	}
	repo.Bot.Send(newMsg)
	return nil
}
