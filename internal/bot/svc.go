package bot

import (
	"context"
	dto2 "github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net/dto"
	"github.com/samber/lo"
	tele "gopkg.in/telebot.v3"
)

func (b *Bot) CreateDevice(ctx context.Context, c tele.Context) error {
	devices, err := b.yandex.Devices(ctx)
	if err != nil {
		return err
	}

	device, ok := lo.Find(devices.Devices, func(d dto2.Device) bool {
		return d.Name == c.Message().Text
	})
	if !ok {
		return c.Send("Не нашли такого устройства")
	}

	if err = b.repo.CreateDevice(ctx, c.Sender().ID, device.ID); err != nil {
		return err
	}

	return b.repo.ChangeState(ctx, c.Sender().ID, 0)
}

func (b *Bot) CreateCommand(ctx context.Context, c tele.Context) error {
	_ = context.Background()
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

	return c.Send("Выберите устройство для сценария", replyMarkup)
}

// 1. добавить устройство
//  -
