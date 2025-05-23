package bot

import (
	"context"
	"fmt"
	"io"

	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro"
	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro/dto"
	yandex_net "github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net"
	"github.com/danzelVash/lampochka/internal/infrastructure/repo"
	"github.com/samber/lo"
	"google.golang.org/appengine/log"

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

func (b *Bot) AddDevice(c tele.Context) error {
	ctx := context.Background()
	devices, err := b.yandex.Devices(ctx)
	if err != nil {
		return err
	}

	var deviceButtons []tele.ReplyButton
	for _, device := range devices.Devices {
		deviceButtons = append(deviceButtons, tele.ReplyButton{Text: device.Name})
	}

	replyMarkup := &tele.ReplyMarkup{
		ReplyKeyboard:   [][]tele.ReplyButton{deviceButtons},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	if err = b.repo.ChangeState(ctx, c.Sender().ID, repo.CreatingDevice); err != nil {
		return err
	}

	return c.Send("Выберите устройство", replyMarkup)
}

func (b *Bot) CreateCommand(c tele.Context) error {
	ctx := context.Background()
	devices, err := b.yandex.Devices(ctx)
	if err != nil {
		return err
	}

	var deviceButtons []tele.ReplyButton
	for _, device := range devices.Devices {
		deviceButtons = append(deviceButtons, tele.ReplyButton{Text: device.Name})
	}

	replyMarkup := &tele.ReplyMarkup{
		ReplyKeyboard:   [][]tele.ReplyButton{deviceButtons},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	if err = b.repo.ChangeState(ctx, c.Sender().ID, repo.CreatingDevice); err != nil {
		return err
	}

	return c.Send("Выберите устройство", replyMarkup)
}

func (b *Bot) OnText(c tele.Context) error {
	ctx := context.Background()
	user, err := b.repo.GetUser(ctx, c.Sender().ID)
	if err != nil {
		return err
	}

	switch user.State {
	case repo.CreatingDevice:
		if err = b.CreateDevice(ctx, c); err != nil {
			return err
		}
		return c.Send("Успешно добавили Ваше устройство")
	}
	return nil
}

func (b *Bot) Start(c tele.Context) error {
	if err := b.repo.CreateUser(context.Background(), c.Sender().ID); err != nil {
		return err
	}
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
