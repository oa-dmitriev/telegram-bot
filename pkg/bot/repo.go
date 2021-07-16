package bot

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	storage "mybot/internal/pkg/storage"

	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotRepo struct {
	Bot       *tgbotapi.BotAPI
	redClient *redis.Client
	db        *sql.DB
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
	conn, _ := storage.InitRedis()
	db, err := storage.InitDB()
	if err != nil {
		return nil, err
	}
	repo := BotRepo{
		Bot:       b,
		redClient: conn,
		db:        db,
	}
	return &repo, nil
}

func (r *BotRepo) CallBack(cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	if cb.Data == "add" {
		return r.Save(cb)
	}
	if cb.Message.ReplyToMessage.Text == "/vocab" {
		return r.CallBackVocab(cb)
	}
	return r.CallBackDef(cb)
}

func (r *BotRepo) Definition(msg *tgbotapi.Message) (tgbotapi.Chattable, error) {
	term := msg.Text
	data, err := r.FetchData(term)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		newMsg := tgbotapi.NewMessage(msg.Chat.ID, "No definition found")
		newMsg.ReplyToMessageID = msg.MessageID
		return newMsg, nil
	}
	dataToSend := GetPage(data, 0)
	newMsg := tgbotapi.NewMessage(
		msg.Chat.ID,
		"•  "+strings.Join(dataToSend, "\n•  "),
	)
	newMsg.ReplyMarkup = CreateMarkup(0, isDataLeft(0, data), true)
	newMsg.ReplyToMessageID = msg.MessageID
	return newMsg, nil
}

func (r *BotRepo) CallBackVocab(cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	num, err := strconv.Atoi(cb.Data)
	if err != nil {
		return nil, err
	}
	data, err := r.GetDataFromDB(cb.From.ID)
	if err != nil {
		return nil, err
	}

	dataToSend := GetPageFromDB(data, num)
	txt := fmt.Sprintf("*%s* - %s", dataToSend[0].Word, dataToSend[0].Definition)
	for i := range dataToSend {
		txt += fmt.Sprintf("\n*%s* - %s", dataToSend[i].Word, dataToSend[i].Definition)
	}
	newMsg := tgbotapi.NewEditMessageText(
		cb.Message.Chat.ID,
		cb.Message.MessageID,
		txt,
	)
	newMsg.ReplyMarkup = CreateMarkup(num, isDataLeft(num, data), false)
	return newMsg, nil
}

func (r *BotRepo) CallBackDef(cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	msg := cb.Message
	term := msg.ReplyToMessage.Text
	num, err := strconv.Atoi(cb.Data)
	if err != nil {
		return nil, err
	}
	data, err := r.FetchData(term)
	if err != nil {
		return nil, err
	}

	dataToSend := GetPage(data, num)
	newMsg := tgbotapi.NewEditMessageText(
		msg.Chat.ID,
		msg.MessageID,
		"•  "+strings.Join(dataToSend, "\n•  "),
	)
	newMsg.ReplyMarkup = CreateMarkup(num, isDataLeft(num, data), true)
	return newMsg, nil
}

func (r *BotRepo) Save(cb *tgbotapi.CallbackQuery) (tgbotapi.Chattable, error) {
	msg := cb.Message
	term := cb.Message.ReplyToMessage.Text
	def, err := r.FetchData(term)
	if err != nil || len(def) == 0 {
		return nil, err
	}
	sqlQuery := `
			SELECT word 
			FROM vocabulary WHERE user_id = $1 AND word = $2
		`
	row := r.db.QueryRow(sqlQuery, cb.From.ID, term)
	var s string
	if err := row.Scan(&s); err == nil {
		newMsg := tgbotapi.NewEditMessageText(
			msg.Chat.ID,
			msg.MessageID,
			fmt.Sprintf("*%s* is already in your vocabulary\n/vocab - show your vocabulary", term),
		)
		newMsg.ParseMode = "markdown"
		log.Println("\n\nALREADY IN DB\n")
		return newMsg, nil
	}
	_, err = r.db.Exec(`
			INSERT INTO vocabulary (user_id, word, definition) 
			VALUES ($1, $2, $3)
		`, cb.From.ID, term, def[0])
	if err != nil {
		return nil, err
	}
	newMsg := tgbotapi.NewEditMessageText(
		msg.Chat.ID,
		msg.MessageID,
		fmt.Sprintf("*%s* added\n/vocab - show your vocabulary", term),
	)
	newMsg.ParseMode = "markdown"
	log.Println("\n\nINSERTED SUCCESSFULY DB\n")
	return newMsg, nil
}

func (r *BotRepo) Command(msg *tgbotapi.Message) (tgbotapi.Chattable, error) {
	if msg.Command() == "vocab" {
		rows, err := r.db.Query(
			"SELECT word FROM vocabulary WHERE user_id = $1",
			msg.From.ID,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		vocab := make([]string, 0)
		for rows.Next() {
			var s string
			err := rows.Scan(&s)
			if err != nil {
				return nil, err
			}
			vocab = append(vocab, s)
		}
		newMsg := tgbotapi.NewMessage(
			msg.Chat.ID,
			"Your vocabulary",
		)
		newMsg.ParseMode = "markdown"

		newMsg.ReplyMarkup = CreateMarkup(0, isDataLeft(0, vocab), false)
		newMsg.ReplyToMessageID = msg.MessageID
		return newMsg, nil
	}
	return nil, fmt.Errorf("Uknown error")
}

func CreateMarkup(curPage int, next, add bool) *tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton
	if curPage == 0 {
		if next {
			keyboard = append(
				keyboard,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Next", strconv.Itoa(curPage+1)),
				),
			)
		}
	} else if next {
		keyboard = append(
			keyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Prev", strconv.Itoa(curPage-1)),
				tgbotapi.NewInlineKeyboardButtonData("Next", strconv.Itoa(curPage+1)),
			),
		)
	} else {
		keyboard = append(
			keyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Prev", strconv.Itoa(curPage-1)),
			),
		)
	}
	if add == true {
		keyboard = append(
			keyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Add", "add"),
			),
		)
	}
	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}
