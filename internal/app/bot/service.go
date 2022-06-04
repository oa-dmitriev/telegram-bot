package botservice

import (
	"context"
	"net/url"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

type Implementation struct {
	userRepo     repository.UserRepo
	vocabRepo    repository.VocabularyRepo
	bot          *tgbotapi.BotAPI
	urbanDictURL *url.URL
	server       *gin.Engine
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
	service := &Implementation{
		userRepo:     userRepo,
		vocabRepo:    vocabRepo,
		bot:          bot,
		urbanDictURL: parsedURL,
		server:       gin.Default(),
	}

	service.server.POST("/", service.Resolve)
	return service, nil
}

func (i *Implementation) Run(ctx context.Context) error {
	return i.server.Run()
}
