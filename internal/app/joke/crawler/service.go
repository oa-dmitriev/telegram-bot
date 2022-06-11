package crawler

import (
	"net/http"
	"net/url"
	"time"

	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

type Implementation struct {
	jokeRepo  repository.JokeRepo
	jokeURL   *url.URL
	client    *http.Client
	rateLimit int
	schedule  time.Duration
}

func NewJokeCrawlerService(
	jokeRepo repository.JokeRepo,
	jokeURL string,
	rateLimit int,
	schedule time.Duration,
) (*Implementation, error) {
	parsedURL, err := url.Parse(jokeURL)
	if err != nil {
		return nil, err
	}

	return &Implementation{
		jokeRepo:  jokeRepo,
		jokeURL:   parsedURL,
		client:    &http.Client{},
		rateLimit: rateLimit,
		schedule:  schedule,
	}, nil
}
