package joke

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/oa-dmitriev/telegram-bot/internal/domain"
	"github.com/oa-dmitriev/telegram-bot/internal/repository"
)

func (i *Implementation) Run(ctx context.Context) error {
	errCount := 0
	for {
		if err := i.StartCrawling(ctx); err != nil {
			errCount++
			log.Printf("error while crawling, tried %d times: %s\n", errCount, err.Error())
		} else {
			errCount = 0
		}

		if errCount > 10 {
			break
		}
		time.Sleep(i.schedule)
	}
	return errors.New("failed crawling 10 times in a row")
}

func (i *Implementation) StartCrawling(ctx context.Context) error {
	for j := 0; j < 1000; j++ {
		jokes, err := i.getResp(ctx)
		if err != nil {
			return err
		}

		for _, joke := range jokes {
			if err := i.jokeRepo.Add(ctx, &repository.DBJoke{
				ID:       joke.ID,
				Type:     joke.Type,
				Joke:     joke.Joke,
				Setup:    joke.Setup,
				Delivery: joke.Delivery,
				Category: joke.Category,
			}); err != nil {
				log.Printf("jokeRepo err: %s\n", err.Error())
			}
		}
		time.Sleep(time.Duration(60 / i.rateLimit))
	}
	return nil
}

func (i *Implementation) getResp(ctx context.Context) ([]*domain.JokeInfo, error) {
	resp, err := http.Get(i.jokeURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jokes := domain.Jokes{}
	if err = json.Unmarshal(body, &jokes); err != nil {
		return nil, err
	}

	return jokes.List, nil
}
