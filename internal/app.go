package internal

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro"
	yandex_net "github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net"
	"github.com/danzelVash/lampochka/internal/infrastructure/repo"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"

	"github.com/danzelVash/lampochka/internal/bot"
	"github.com/jackc/pgx/v5"
	tele "gopkg.in/telebot.v3"

	googlegrpc "google.golang.org/grpc"
)

func init() {
	// ignore db columns that doesn't exist at the destination
	dbscanAPI, err := pgxscan.NewDBScanAPI(dbscan.WithAllowUnknownColumns(true))
	if err != nil {
		panic(err)
	}

	api, err := pgxscan.NewAPI(dbscanAPI)
	if err != nil {
		panic(err)
	}

	pgxscan.DefaultAPI = api
}

type App struct {
	BotSvc *bot.Bot
	TgBot  *tele.Bot
	Repo   *repo.Repo

	neuroConn *googlegrpc.ClientConn
	pgxConn   *pgx.Conn

	Neuro     *neuro.Gateway
	YandexNet *yandex_net.Gateway
}

func NewApp(ctx context.Context) *App {
	fmt.Println("Инициализация бота")
	pref := tele.Settings{
		Token:  "7765937182:AAFRkUKr3iUdxfYWiTJdxeLNMB-5cuF_M6g",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	botTg, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Инициализация постгреса")
	// pgx
	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@postgres:5432/mirea?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Инициализация клиента НЕЙРО")
	// neuro gate
	neuroConn, err := googlegrpc.NewClient("51.250.93.99:8000",
		googlegrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		TgBot:     botTg,
		neuroConn: neuroConn,
		pgxConn:   conn,
	}
}

func (a *App) Init() {
	fmt.Println("ИНИТ")

	a.Repo = repo.New(a.pgxConn)
	a.YandexNet = yandex_net.NewGateway()
	a.Neuro = neuro.NewGateway(neuro.NewExternalClient(a.neuroConn))
	a.BotSvc = bot.New(a.TgBot, a.Neuro, a.YandexNet, a.Repo)
}
