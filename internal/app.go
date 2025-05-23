package internal

import (
	"context"
	"net/http"
)

type App struct {
	server *http.Server
}

func NewApp(ctx context.Context) *App {
	app := &App{}

	return app
}

func (app *App) Run(ctx context.Context) error {
	return app.server.ListenAndServe()
}
