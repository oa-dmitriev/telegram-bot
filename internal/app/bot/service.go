package botservice

import (
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

type Implementation struct {
	userRepo  repository.UserRepo
	vocabRepo repository.VocabularyRepo
	bot       *tgbotapi.BotAPI

	urbanDictURL *url.URL
}

func NewBotService(
	userRepo repository.UserRepo,
	vocabRepo repository.VocabularyRepo,
	bot *tgbotapi.BotAPI,
	urbanDictURL string,
) (*Implementation, error) {
	parsedURL, err := url.Parse(urbanDictURL)
	if err != nil {
		return nil, err
	}
	return &Implementation{
		userRepo:     userRepo,
		vocabRepo:    vocabRepo,
		bot:          bot,
		urbanDictURL: parsedURL,
	}, nil
}
