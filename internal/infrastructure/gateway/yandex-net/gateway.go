package yandex_net

import (
	"context"
	"encoding/json"
	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net/dto"
	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/yandex-net/dto/on_off"
	"github.com/go-resty/resty/v2"
)

const (
	token = "Bearer y0__xC58d2EARjw9jcgyci_mxMlSyMEljQVGTpgpQAG2zJlECWgRQ"

	onOffEndpoint = "https://api.iot.yandex.net/v1.0/devices/actions"

	devicesEndpoint = "https://api.iot.yandex.net/v1.0/user/info"
)

type Gateway struct {
	client *resty.Client
}

func NewGateway() *Gateway {
	return &Gateway{client: resty.New()}
}

func (g Gateway) OnOffDevice(ctx context.Context, id string, value bool) error {
	_, err := g.client.R().SetContext(ctx).SetHeader("Authorization", token).
		SetBody(on_off.New(id, value)).
		Post(onOffEndpoint)
	return err
}

func (g Gateway) Devices(ctx context.Context) (devices dto.Devices, err error) {
	response, err := g.client.R().SetContext(ctx).SetHeader("Authorization", token).
		Get(devicesEndpoint)
	if err != nil {
		return dto.Devices{}, err
	}
	return devices, json.Unmarshal(response.Body(), &devices)
}
