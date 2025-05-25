package main

import (
	"context"
	"fmt"

	"github.com/danzelVash/lampochka/internal"
	tele "gopkg.in/telebot.v3"
)

func main() {
	ctx := context.Background()
	fmt.Println("Старт сервера")
	app := internal.NewApp(ctx)
	app.Init()

	app.TgBot.Handle("/start", app.BotSvc.Start)

	app.TgBot.Handle("/help", app.BotSvc.Help)

	app.TgBot.Handle("/exit", app.BotSvc.Help)

	app.TgBot.Handle("/adddevice", app.BotSvc.AddDevice)

	app.TgBot.Handle("/createcommand", app.BotSvc.CreateCommand)

	app.TgBot.Handle("/mycommands", app.BotSvc.MyCommands)

	app.TgBot.Handle("/deletecommand", app.BotSvc.DeleteCommand)

	app.TgBot.Handle(tele.OnVoice, app.BotSvc.VoiceMess)

	app.TgBot.Handle(tele.OnText, app.BotSvc.OnText)

	app.TgBot.Start()
}
