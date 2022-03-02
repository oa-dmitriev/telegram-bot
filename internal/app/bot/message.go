package botservice

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oa-dmitriev/telegram-bot/internal/domain"
	"github.com/oa-dmitriev/telegram-bot/internal/markup"
)

const pageLen = 3

func (i *Implementation) Definition(ctx context.Context, msg *tgbotapi.Message, pageNum int) (tgbotapi.Chattable, error) {
	term := msg.Text
	data, err := i.FetchData(ctx, term)
	if err != nil {
		return nil, err
	}

	sendMsg := tgbotapi.NewMessage(msg.Chat.ID, "")
	sendMsg.ReplyToMessageID = msg.MessageID

	if len(data) == 0 {
		sendMsg.Text = "No definition found"
		return sendMsg, nil
	}

	dataToSend := getPage(data, pageLen, pageNum)
	sendMsg.Text = "•  " + strings.Join(dataToSend, "\n•  ")

	mark := markup.New().WithPrev(pageNum)
	if len(data) > len(dataToSend) {
		mark = mark.WithNext(0)
	}
	mark = mark.WithCustomMsg("Add")

	sendMsg.ReplyMarkup = mark.InlineKeyboardMarkup
	return sendMsg, nil
}

func getPage(data []string, limit, pageNum int) []string {
	start := pageNum * pageLen
	if start > len(data) {
		return nil
	}
	end := (pageNum + 1) * pageLen
	if len(data) < end {
		end = len(data)
	}
	return data[start:end]
}

func (i *Implementation) FetchData(ctx context.Context, word string) ([]string, error) {
	u := i.urbanDictURL
	p := url.Values{
		"term": []string{word},
	}
	u.RawQuery = p.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := domain.Data{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	defs := data.Definitions

	sort.Slice(defs, func(i, j int) bool {
		return defs[i].ThumbsUp > defs[j].ThumbsUp
	})

	res := make([]string, len(data.Definitions))
	for i := range data.Definitions {
		res[i] = data.Definitions[i].Definition
	}

	// TODO add res to cache

	return res, nil
}
