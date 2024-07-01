package app

import (
	"github.com/zhikh23/itsreg-bots/internal/app/grpcapp"
	"github.com/zhikh23/itsreg-bots/internal/services/processor"
	"log/slog"
)

type App struct {
	GrpcApp *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
) *App {
	service, err := processor.New(
		processor.WithLogger(log),
		processor.WithMemoryBotRepository(),
		processor.WithMemoryModuleRepository(),
		processor.WithMemoryModuleRepository())

	if err != nil {
		panic(err)
	}

	grpcApp := grpcapp.New(log, service, grpcPort)

	return &App{
		GrpcApp: grpcApp,
	}
}
