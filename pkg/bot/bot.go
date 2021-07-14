package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	apiURL = url.URL{
		Scheme: "https",
		Host:   "api.urbandictionary.com",
		Path:   "v0/define",
	}
	PAGELEN = 3
)

type DefinitionData struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type Data struct {
	Definitions []*DefinitionData `json:"list"`
}

var cmds = map[string]func(*BotRepo, *tgbotapi.Message, []string) error{
	"/def": func(repo *BotRepo, msgInfo *tgbotapi.Message, args []string) error {
		data, err := FetchData(strings.Join(args, " "))
		if err != nil {
			return fmt.Errorf("server error: %#v\n", err)
		}
		if len(data) == 0 {
			msg := tgbotapi.NewMessage(msgInfo.Chat.ID, "No definition found")
			msg.ReplyToMessageID = msgInfo.MessageID
			repo.Bot.Send(msg)
			return nil
		}
		dataToSend := GetPage(data, 0)
		msg := tgbotapi.NewMessage(
			msgInfo.Chat.ID,
			"•  "+strings.Join(dataToSend, "\n•  "),
		)
		if isDataLeft(data, 0) {
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Next", "1"),
				),
			)
		}
		msg.ReplyToMessageID = msgInfo.MessageID
		repo.Bot.Send(msg)
		return nil
	},
	"/add": func(repo *BotRepo, msg *tgbotapi.Message, args []string) error {

		return nil
	},
}

func isDataLeft(data []string, offset int) bool {
	return len(data) > offset+PAGELEN
}

func FetchData(term string) ([]string, error) {
	u := apiURL
	p := url.Values{
		"term": []string{term},
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
	data := Data{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	ans := make([]string, len(data.Definitions))
	for i := range data.Definitions {
		ans[i] = data.Definitions[i].Definition
	}
	return ans, nil
}

func GetPage(data []string, offset int) []string {
	start := offset * PAGELEN
	end := offset*PAGELEN + PAGELEN
	if len(data) < end {
		end = len(data)
	}
	return data[start:end]
}
