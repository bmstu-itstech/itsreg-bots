package grpc

import (
	"context"
	"github.com/zhikh23/itsreg-bots/internal/objects"
	botsproto "github.com/zhikh23/itsreg-proto/gen/go/bots/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Processor interface {
	Process(ctx context.Context, botId int64, userId int64, ans string) ([]objects.Message, error)
}

type serverAPI struct {
	botsproto.UnimplementedBotsServiceServer
	p Processor
}

func Register(server *grpc.Server, processor Processor) {
	botsproto.RegisterBotsServiceServer(server, &serverAPI{
		p: processor,
	})
}

func (s *serverAPI) Process(ctx context.Context, req *botsproto.ProcessRequest) (*botsproto.ProcessResponse, error) {
	messages, err := s.p.Process(ctx, req.BotId, req.UserId, req.Text)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &botsproto.ProcessResponse{
		Messages: messagesToDto(messages),
	}

	return res, nil
}

func messagesToDto(messages []objects.Message) []*botsproto.Message {
	dtos := make([]*botsproto.Message, len(messages))

	for i, message := range messages {
		dtos[i] = messageToDto(&message)
	}

	return dtos
}

func messageToDto(message *objects.Message) *botsproto.Message {
	buttons := make([]*botsproto.Button, len(message.Buttons))
	for i, btn := range message.Buttons {
		buttons[i] = buttonToDto(&btn)
	}

	return &botsproto.Message{
		Text:    message.Text,
		Buttons: buttons,
	}
}

func buttonToDto(button *objects.Button) *botsproto.Button {
	return &botsproto.Button{
		Text: button.Text,
	}
}
