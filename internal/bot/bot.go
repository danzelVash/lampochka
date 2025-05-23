package bot

import (
	"context"
	"fmt"
	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro"
	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro/dto"
	yandex_net "github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net"
	"github.com/danzelVash/lampochka/internal/infrastructure/repo"
	"github.com/samber/lo"
	"google.golang.org/appengine/log"
	"io"

	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	tgBot *tele.Bot

	neuro  *neuro.Gateway
	yandex *yandex_net.Gateway
	repo   *repo.Repo
}

func New(tgBot *tele.Bot, neuro *neuro.Gateway, yandex *yandex_net.Gateway, repo *repo.Repo) *Bot {
	return &Bot{tgBot: tgBot, neuro: neuro, yandex: yandex, repo: repo}
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

	commands, err := b.repo.GetCommands(context.Background(), c.Chat().ID)
	if err != nil {
		return err
	}

	matched, err := b.neuro.GetAudio(context.Background(), lo.Map(commands, func(command repo.Command, _ int) dto.Command {
		return dto.Command{Name: command.Command}
	}), bytes)

	log.Infof(context.Background(), fmt.Sprintf("Подобранный сценарий: %s", matched))

	return b.yandex.Match(context.Background(), lo.FindOrElse(commands, repo.Command{}, func(command repo.Command) bool {
		return command.Command == matched.Name
	}))
}

func (b *Bot) Create(c tele.Context) error {
	//// Создаем меню с устройствами
	//var deviceButtons []tele.ReplyButton
	//for _, device := range devices {
	//	deviceButtons = append(deviceButtons, tele.ReplyButton{Text: device})
	//}
	//
	//replyMarkup := &tele.ReplyMarkup{
	//	ReplyKeyboard:   [][]tele.ReplyButton{deviceButtons},
	//	ResizeKeyboard:  true,
	//	OneTimeKeyboard: true,
	//}
	//return c.Send(helpText)
	return nil
}

func (b *Bot) OnText(c tele.Context) error {
	// todo: switch по статусу
	return nil
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
