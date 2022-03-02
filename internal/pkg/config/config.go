package config

import (
	"context"

	"github.com/kelseyhightower/envconfig"
)

const (
	envPrefix = "PERIC"
)

type Env struct {
	BotToken     string `envconfig:"BOT_TOKEN"`
	WebHookURL   string `envconfig:"WEB_HOOK_URL"`
	DatabaseURL  string `envconfig:"DATABASE_URL"`
	UrbanDictURL string `envconfig:"URBAN_DICT_URL" default:"https://api.urbandictionary.com/v0/define"`
}

func GetEnv(ctx context.Context) (*Env, error) {
	cfg := Env{}
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
