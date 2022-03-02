package main

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oa-dmitriev/telegram-bot/internal/app"
	botservice "github.com/oa-dmitriev/telegram-bot/internal/app/bot"
	"github.com/oa-dmitriev/telegram-bot/internal/pkg/config"
	"github.com/oa-dmitriev/telegram-bot/internal/pkg/database"
	"github.com/oa-dmitriev/telegram-bot/internal/repository/user"
	"github.com/oa-dmitriev/telegram-bot/internal/repository/vocabulary"
)

func main() {
	a := app.New()
	ctx := context.Background()

	cfg, err := config.GetEnv(ctx)
	exitOnError(ctx, "could not read config: %s", err)

	log.Println("config: ", cfg)

	db, err := database.NewDBWrapper(ctx, database.Options{
		DatabaseURL: cfg.DatabaseURL,
	})
	exitOnError(ctx, "could not create connection to DB: %s", err)

	exitOnError(ctx, "could not connect to DB: %s", database.Ready(db))

	userRepo := user.NewRepo(db)
	vocabRepo := vocabulary.NewRepo(db)

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	exitOnError(ctx, "could not create BotAPI: %s", err)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(cfg.WebHookURL))
	exitOnError(ctx, "could not set webhook: %s", err)

	botService, err := botservice.NewBotService(userRepo, vocabRepo, bot, cfg.UrbanDictURL)
	exitOnError(ctx, "could not create botService: %s", err)

	a.Server().POST("/", botService.Resolve)
	exitOnError(ctx, "could not run the app: %s", a.Run())
}

func exitOnError(_ context.Context, msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err)
	}
}
