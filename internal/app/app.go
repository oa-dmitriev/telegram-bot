package app

import (
	"context"

	"github.com/gin-gonic/gin"
)

type App struct {
	publicServer *gin.Engine
}

func New() *App {
	return &App{
		publicServer: gin.Default(),
	}
}

func (a *App) Server() *gin.Engine {
	return a.publicServer
}

func (a *App) Run(ctx context.Context, services ...Service) error {
	errChan := make(chan error)
	for _, service := range services {
		go func(service Service) {
			if err := service.Run(ctx); err != nil {
				errChan <- err
			}
		}(service)
	}

	err := <-errChan
	return err
}
