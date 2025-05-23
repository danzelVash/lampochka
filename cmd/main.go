package main

import (
	"context"

	"github.com/danzelVash/lampochka/internal"
	tele "gopkg.in/telebot.v3"
)

func main() {
	ctx := context.Background()
	app := internal.NewApp(ctx)
	app.Init()

	app.TgBot.Handle("/start", app.BotSvc.Start)

	app.TgBot.Handle("/help", app.BotSvc.Help)

	app.TgBot.Handle("/addDevice", app.BotSvc.AddDevice)

	app.TgBot.Handle("/createCommand", app.BotSvc.CreateCommand)

	app.TgBot.Handle("/addCommand", app.BotSvc.CreateCommand)
	app.TgBot.Handle(tele.OnVoice, app.BotSvc.VoiceMess)

	app.TgBot.Handle(tele.OnText, app.BotSvc.OnText)

	app.TgBot.Start()
}
