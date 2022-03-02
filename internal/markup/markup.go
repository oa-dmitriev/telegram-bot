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
		} else {
			return nil
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
	if add {
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
