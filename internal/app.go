package internal

import (
	"context"
	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro"
	yandex_net "github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net"
	"github.com/danzelVash/lampochka/internal/infrastructure/repo"
	"log"
	"time"

	"github.com/danzelVash/lampochka/internal/bot"
	"github.com/jackc/pgx/v5"
	tele "gopkg.in/telebot.v3"

	googlegrpc "google.golang.org/grpc"
)

type App struct {
	BotSvc *bot.Bot
	TgBot  *tele.Bot
	Repo   *repo.Repo

	Neuro     *neuro.Gateway
	YandexNet *yandex_net.Gateway
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

	// neuro gate
	neuroConn, err := googlegrpc.NewClient("localhost:7002")
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		BotSvc:    bot.New(botTg),
		TgBot:     botTg,
		Repo:      repo.New(conn),
		YandexNet: yandex_net.NewGateway(),
		Neuro:     neuro.NewGateway(neuro.NewExternalClient(neuroConn)),
	}
}
