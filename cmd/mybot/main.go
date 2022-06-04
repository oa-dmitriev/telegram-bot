package main

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oa-dmitriev/telegram-bot/internal/app"
	botservice "github.com/oa-dmitriev/telegram-bot/internal/app/bot"
	crawler "github.com/oa-dmitriev/telegram-bot/internal/app/joke"
	"github.com/oa-dmitriev/telegram-bot/internal/pkg/config"
	"github.com/oa-dmitriev/telegram-bot/internal/pkg/database"
	"github.com/oa-dmitriev/telegram-bot/internal/repository/joke"
	"github.com/oa-dmitriev/telegram-bot/internal/repository/user"
	"github.com/oa-dmitriev/telegram-bot/internal/repository/vocabulary"
)

func main() {
	a := app.New()
	ctx := context.Background()

	cfg, err := config.GetEnv(ctx)
	exitOnError(ctx, "could not read config: %s", err)

	db, err := database.NewDBWrapper(ctx, database.Options{
		DatabaseURL: cfg.DatabaseURL,
	})
	exitOnError(ctx, "could not create connection to DB: %s", err)

	exitOnError(ctx, "could not connect to DB: %s", database.Ready(db))

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	exitOnError(ctx, "could not create BotAPI: %s", err)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(cfg.WebHookURL))
	exitOnError(ctx, "could not set webhook: %s", err)

	//-------------------- init services -----------------------
	userRepo := user.NewRepo(db)
	vocabRepo := vocabulary.NewRepo(db)
	jokeRepo := joke.NewRepo(db)

	jokeCrawlerService, err := crawler.NewJokeCrawlerService(
		jokeRepo,
		cfg.JokeAPI,
		cfg.JokeAPIRateLimit,
	)
	exitOnError(ctx, "could not create jokeCrawlerService: %s", err)

	botService, err := botservice.NewBotService(userRepo, vocabRepo, bot, cfg.UrbanDictURL)
	exitOnError(ctx, "could not create botService: %s", err)

	err = a.Run(ctx, jokeCrawlerService, botService)
	exitOnError(ctx, "could not run the app: %s", err)
}

func exitOnError(_ context.Context, msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
