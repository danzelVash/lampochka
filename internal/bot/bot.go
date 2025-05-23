package bot

import (
	"fmt"
	"io"

	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	tgBot *tele.Bot
}

func New(tgBot *tele.Bot) *Bot {
	return &Bot{tgBot: tgBot}
}

func (b *Bot) VoiceMess(c tele.Context) error {
	voice := c.Message().Voice

	file, err := b.tgBot.FileByID(voice.FileID)
	if err != nil {
		return c.Send("Не удалось получить файл голосового сообщения")
	}

	f, err := b.tgBot.File(&file)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
	return err
}

func (b *Bot) Start(c tele.Context) error {
	return c.Send("Привет! Я пример бота на Go. Отправь мне /help для списка команд.")
}

func (b *Bot) Help(c tele.Context) error {
	helpText := `Доступные команды:
/start - начать работу с ботом
/help - показать это сообщение
/echo [текст] - повторить текст
/time - показать текущее время`
	return c.Send(helpText)
}
