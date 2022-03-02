package app

import "github.com/gin-gonic/gin"

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

func (a *App) Run() error {
	return a.publicServer.Run()
}
