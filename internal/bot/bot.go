package bot

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro"
	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro/dto"
	yandex_net "github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net"
	"github.com/danzelVash/lampochka/internal/infrastructure/repo"
	"github.com/samber/lo"
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
	fmt.Println("дернулась ручка VoiceMess")
	ctx := context.Background()
	commands, err := b.repo.GetCommands(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	if len(commands) == 0 {
		return c.Send("У вас нет ни одного сценария")
	}

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

	matched, err := b.neuro.GetAudio(ctx, lo.Map(commands, func(command repo.Command, _ int) dto.Command {
		return dto.Command{Name: command.Command}
	}), bytes)
	if err != nil {
		return err
	}

	err = b.yandex.Match(ctx, lo.FindOrElse(commands, repo.Command{}, func(command repo.Command) bool {
		return command.Command == matched.Name
	}))
	if err != nil {
		return err
	}

	return c.Send(fmt.Sprintf("Выполнена команда \"%s\"", matched.Name))
}

func (b *Bot) AddDevice(c tele.Context) error {
	ctx := context.Background()
	fmt.Println("дернулась ручка AddDevice")

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
	fmt.Println("дернулась ручка CreateCommand")

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

	if err = b.repo.ChangeState(ctx, c.Sender().ID, repo.CreatingCommandDevice); err != nil {
		return err
	}

	return c.Send("Выберите устройство для сценария", replyMarkup)
}

func (b *Bot) MyCommands(c tele.Context) error {
	ctx := context.Background()
	fmt.Println("дернулась ручка MyCommands")

	commands, err := b.repo.GetCommands(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	if len(commands) == 0 {
		return c.Send("У вас нет ни одного сценария")
	}

	var message strings.Builder
	message.WriteString("Список ваших сценариев:\n\n")

	for _, command := range commands {
		message.WriteString(fmt.Sprintf("%s - %s\n", command.Command, command.Action))
	}

	return c.Send(message.String())
}

func (b *Bot) DeleteCommand(c tele.Context) error {
	ctx := context.Background()
	fmt.Println("дернулась ручка MyCommands")

	if err := b.repo.ChangeState(ctx, c.Sender().ID, repo.DeleteCommand); err != nil {
		return err
	}

	return c.Send("Введите название сценария, который хотите удалить")
}

func (b *Bot) CreateAction(c tele.Context) error {
	ctx := context.Background()
	fmt.Println("дернулась ручка CreateAction")

	actions, err := b.repo.GetCommandList(ctx)
	if err != nil {
		return err
	}

	var actionButtons []tele.ReplyButton
	for _, action := range actions {
		actionButtons = append(actionButtons, tele.ReplyButton{Text: action.Action})
	}

	replyMarkup := &tele.ReplyMarkup{
		ReplyKeyboard:   [][]tele.ReplyButton{actionButtons},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}

	return c.Send("Выберите сценарий", replyMarkup)
}

func (b *Bot) OnText(c tele.Context) error {
	ctx := context.Background()
	fmt.Println("дернулась ручка OnText")

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
	case repo.CreatingCommandDevice:
		if err = b.CreateCommandDevice(ctx, c); err != nil {
			return err
		}
		return b.CreateAction(c)
	case repo.CreatingCommandAction:
		if err = b.CreateCommandAction(ctx, c); err != nil {
			return err
		}
		return c.Send("Как хотите вызывать команду?")
	case repo.CreatingCommandText:
		if err = b.CreateCommandText(ctx, c); err != nil {
			return err
		}
		return c.Send("Ваш сценарий успешно заведен")
	case repo.DeleteCommand:
		if err = b.DeleteCommandService(ctx, c); err != nil {
			return c.Send("Произошла ошибка при удалении сценария")
		}
		return c.Send("Ваш сценарий успешно удален")
	}
	return nil
}

func (b *Bot) Start(c tele.Context) error {
	fmt.Println("дернулась ручка Start")

	if err := b.repo.CreateUser(context.Background(), c.Sender().ID); err != nil {
		return err
	}
	return c.Send("Привет! Я бот, созданный Максимом Нечепоруком, Даней Узяновым и Даней Булыкиным для МИРЭА")
}

func (b *Bot) Exit(c tele.Context) error {
	fmt.Println("дернулась ручка Exit")

	err := b.repo.ChangeState(context.Background(), c.Sender().ID, 0)
	if err != nil {
		return err
	}
	return c.Send("Очистил ваше состояние")
}

func (b *Bot) Help(c tele.Context) error {
	fmt.Println("дернулась ручка Help")

	helpText := `Доступные команды:
/start - Начать работу с ботом
/adddevice - Добавить устройство
/createdommand - Добавить сценарий умного устройства
/exit - Очистка состояний
/mycommands - Мои сценарии
/deletecommand - Удалить сценарий`
	return c.Send(helpText)
}
