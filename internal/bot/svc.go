package bot

import (
	"context"

	dto2 "github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net/dto"
	"github.com/danzelVash/lampochka/internal/infrastructure/repo"
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

func (b *Bot) CreateCommandDevice(ctx context.Context, c tele.Context) error {
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

	if err = b.repo.CreateCommandDevice(ctx, c.Sender().ID, device.ID); err != nil {
		return err
	}

	return b.repo.ChangeState(ctx, c.Sender().ID, repo.CreatingCommandAction)
}

func (b *Bot) CreateCommandAction(ctx context.Context, c tele.Context) error {
	actions, err := b.repo.GetCommandList(ctx)
	if err != nil {
		return err
	}

	_, ok := lo.Find(actions, func(act repo.Command) bool {
		return act.Action == c.Message().Text
	})
	if !ok {
		return c.Send("Не нашли такого сценария")
	}

	if err = b.repo.CreateCommandAction(ctx, c.Sender().ID, c.Message().Text); err != nil {
		return err
	}

	return b.repo.ChangeState(ctx, c.Sender().ID, repo.CreatingCommandText)
}

func (b *Bot) CreateCommandText(ctx context.Context, c tele.Context) error {
	if err := b.repo.CreateCommandText(ctx, c.Sender().ID, c.Message().Text); err != nil {
		return err
	}

	return b.repo.ChangeState(ctx, c.Sender().ID, 0)
}

// 1. добавить устройство
//  -
