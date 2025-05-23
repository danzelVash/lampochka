package main

import (
	"log"
	"time"

	"github.com/danzelVash/lampochka/internal/bot"
	tele "gopkg.in/telebot.v3"
)

func main() {
	pref := tele.Settings{
		Token:  "7765937182:AAFRkUKr3iUdxfYWiTJdxeLNMB-5cuF_M6g",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	botTg, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	botSvc := bot.New(botTg)

	botTg.Handle("/start", botSvc.Start)

	// Обработчик команды /help
	botTg.Handle("/help", botSvc.Help)

	// Обработчик любого аудио сообщения
	botTg.Handle(tele.OnVoice, botSvc.VoiceMess)

	// Обработчик любого текстового сообщения
	botTg.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Я не понимаю эту команду. Попробуйте /help")
	})

	log.Println("Бот запущен...")
	botTg.Start()
}
