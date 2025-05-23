package neuro

import (
	"context"
	"fmt"
	"github.com/danzelVash/lampochka/internal/infrastructure/gateway/neuro/dto"
	"github.com/danzelVash/lampochka/proto/pb/github.com/danzelVash/lampochka/proto"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

func NewExternalClient(conn grpc.ClientConnInterface) proto.AudioRecognizerClient {
	return proto.NewAudioRecognizerClient(conn)
}

type Gateway struct {
	client proto.AudioRecognizerClient
}

func NewGateway(client proto.AudioRecognizerClient) *Gateway {
	return &Gateway{client: client}
}

func (g Gateway) GetAudio(ctx context.Context, commands []dto.Command, audio []byte) (dto.Command, error) {
	fmt.Printf("[NEOURO GATEWAY] GetAudio, commands:%v\n", commands)

	response, err := g.client.GetAudio(ctx, &proto.GetAudioRequest{
		Chunk: audio,
		Commands: lo.Map(commands, func(command dto.Command, _ int) *proto.GetAudioRequest_Command {
			return &proto.GetAudioRequest_Command{Name: command.Name}
		}),
	})
	if err != nil {
		return dto.Command{}, err
	}

	return dto.Command{Name: response.GetCommand()}, nil
}
