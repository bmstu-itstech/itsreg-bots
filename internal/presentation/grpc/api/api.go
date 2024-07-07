package api

import (
	"context"
	"github.com/zhikh23/itsreg-bots/internal/application"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	botsv1 "github.com/zhikh23/itsreg-proto/gen/go/bots"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcApi struct {
	botsv1.UnimplementedBotsServiceServer
	bots *application.BotsService
}

func Register(
	grpcServer *grpc.Server,
	bots *application.BotsService,
) {
	botsv1.RegisterBotsServiceServer(grpcServer, &grpcApi{bots: bots})
}

func (a *grpcApi) Create(ctx context.Context, req *botsv1.CreateRequest) (*botsv1.CreateResponse, error) {
	blocks := blocksToDtos(req.Blocks)

	id, err := a.bots.Create(ctx, req.Name, req.Token, req.Start, blocks)
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

	messages, err := a.bots.Process(ctx, uint64(botId), uint64(userId), req.Text)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process bot: %v", err)
	}

	return &botsv1.ProcessResponse{
		Messages: messagesFromDtos(messages),
	}, nil
}
