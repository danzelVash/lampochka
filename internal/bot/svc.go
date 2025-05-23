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

	return b.repo.CreateDevice(ctx, c.Sender().ID, device.ID)
}
