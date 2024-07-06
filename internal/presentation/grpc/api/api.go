package api

import (
	"context"
	"github.com/zhikh23/itsreg-bots/internal/application"
	botsv1 "github.com/zhikh23/itsreg-proto/gen/go/bots/proto"
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

func (g *grpcApi) Process(ctx context.Context, req *botsv1.ProcessRequest) (*botsv1.ProcessResponse, error) {
	messages, err := g.bots.Process(ctx, req.BotId, req.UserId, req.Text)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &botsv1.ProcessResponse{
		Messages: messagesToDtos(messages),
	}, nil
}
