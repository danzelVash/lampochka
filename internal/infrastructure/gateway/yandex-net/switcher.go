package yandex_net

import (
	"context"
	"github.com/danzelVash/lampochka/internal/infrastructure/repo"
)

func (g Gateway) Match(ctx context.Context, command repo.Command) error {
	switch command.Command {
	case "Включить":
		return g.OnOffDevice(ctx, command.Device, true)
	case "Выключить":
		return g.OnOffDevice(ctx, command.Device, false)
	}
	return nil
}
