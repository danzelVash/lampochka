package main

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

func main() {
	// Получаем токен бота из переменной окружения
	token := "7765937182:AAFRkUKr3iUdxfYWiTJdxeLNMB-5cuF_M6g"

	// Настройки бота
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	// Создаем экземпляр бота
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Обработчик команды /start
	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! Я пример бота на Go. Отправь мне /help для списка команд.")
	})

	// Обработчик команды /help
	bot.Handle("/help", func(c tele.Context) error {
		helpText := `Доступные команды:
		/start - начать работу с ботом
		/help - показать это сообщение
		/echo [текст] - повторить текст
		/time - показать текущее время`
		return c.Send(helpText)
	})

	// Обработчик любого текстового сообщения
	bot.Handle(tele.OnText, func(c tele.Context) error {
		return c.Send("Я не понимаю эту команду. Попробуйте /help")
	})

	log.Println("Бот запущен...")
	bot.Start()
}
