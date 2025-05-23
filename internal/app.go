package internal

import (
	"context"
	"log"
	"time"

	"github.com/danzelVash/lampochka/internal/bot"
	"github.com/danzelVash/lampochka/internal/repo"
	"github.com/jackc/pgx/v5"
	tele "gopkg.in/telebot.v3"
)

type App struct {
	BotSvc *bot.Bot
	TgBot  *tele.Bot
	Repo   *repo.Repo
}

func NewApp(ctx context.Context) *App {
	pref := tele.Settings{
		Token:  "7765937182:AAFRkUKr3iUdxfYWiTJdxeLNMB-5cuF_M6g",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	botTg, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/mirea?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		BotSvc: bot.New(botTg),
		TgBot:  botTg,
		Repo:   repo.New(conn),
	}

	return app
}
