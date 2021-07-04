package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var dataUrl = "https://wordsapiv1.p.mashape.com/words"

type BotRepo struct {
	Bot *tgbotapi.BotAPI
}

type DefinitionData struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

type Data struct {
	Definitions []*DefinitionData `json:"list"`
}

var cmds = map[string]func(*BotRepo, *tgbotapi.Message, []string) error{
	"/def": func(repo *BotRepo, msgInfo *tgbotapi.Message, args []string) error {
		log.Println("/def: ", args)
		p := url.Values{
			"term": []string{strings.Join(args, " ")},
		}
		u := url.URL{
			Scheme:   "https",
			Host:     "api.urbandictionary.com",
			Path:     "v0/define",
			RawQuery: p.Encode(),
		}
		resp, err := http.Get(u.String())
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		data := Data{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		for _, v := range data.Definitions {
			log.Println(v.Definition)
		}

		msg := tgbotapi.NewMessage(msgInfo.Chat.ID, data.Definitions[0].Definition)
		msg.ReplyToMessageID = msgInfo.MessageID
		repo.Bot.Send(msg)
		return nil
	},
	"/add": func(repo *BotRepo, msg *tgbotapi.Message, args []string) error {

		return nil
	},
}
