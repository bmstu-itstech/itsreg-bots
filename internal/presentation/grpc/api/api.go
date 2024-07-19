package api

import (
	"context"
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/application"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	botsv1 "github.com/bmstu-itstech/itsreg-proto/gen/go/bots"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcApi struct {
	botsv1.UnimplementedBotsServiceServer
	app *application.App
}

func Register(
	grpcServer *grpc.Server,
	app *application.App,
) {
	botsv1.RegisterBotsServiceServer(grpcServer, &grpcApi{app: app})
}

func (a *grpcApi) Create(ctx context.Context, req *botsv1.CreateRequest) (*botsv1.CreateResponse, error) {
	blocks := blocksToDtos(req.Blocks)

	id, err := a.app.Registry.Create(ctx, req.Name, req.Token, req.Start, blocks)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create bot: %v", err)
	}

	return &botsv1.CreateResponse{
		BotId: id,
	}, nil
}

func (a *grpcApi) Process(ctx context.Context, req *botsv1.ProcessRequest) (*botsv1.ProcessResponse, error) {
	botId := value.BotId(req.BotId)
	userId := value.UserId(req.UserId)

	messages, err := a.app.Processor.Process(ctx, uint64(botId), uint64(userId), req.Text)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process bot: %v", err)
	}

	return &botsv1.ProcessResponse{
		Messages: messagesFromDtos(messages),
	}, nil
}

func (a *grpcApi) GetToken(ctx context.Context, req *botsv1.TokenRequest) (*botsv1.TokenResponse, error) {
	botId := value.BotId(req.BotId)

	token, err := a.app.Registry.Token(ctx, uint64(botId))
	if err != nil {
		if errors.Is(err, interfaces.ErrBotNotFound) {
			return nil, status.Errorf(codes.NotFound, "bot not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get token: %v", err)
	}

	return &botsv1.TokenResponse{
		Token: token,
	}, err
}
