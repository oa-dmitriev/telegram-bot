package markup

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Markup struct {
	*tgbotapi.InlineKeyboardMarkup
}

func New() *Markup {
	return &Markup{
		&tgbotapi.InlineKeyboardMarkup{},
	}
}

func (m *Markup) WithPrev(curPage int) *Markup {
	if curPage <= 0 {
		return m
	}
	m.InlineKeyboard = append(m.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Prev", strconv.Itoa(curPage-1)),
	))
	return m
}

func (m *Markup) WithNext(curPage int) *Markup {
	m.InlineKeyboard = append(m.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Next", strconv.Itoa(curPage+1)),
	))
	return m
}

func (m *Markup) WithCustomMsg(customMsg string) *Markup {
	m.InlineKeyboard = append(m.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(customMsg, strings.ToLower(customMsg)),
	))
	return m
}
