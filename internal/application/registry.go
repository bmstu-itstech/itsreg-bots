package application

import (
	"github.com/zhikh23/itsreg-bots/internal/config"
	"github.com/zhikh23/itsreg-bots/internal/domain/service"
	ansmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/answer/memory"
	blockmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/block/memory"
	botmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/bot/memory"
	prtmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/participant/memory"
	"log/slog"
)

var (
	Logger         *slog.Logger
	BotsAppService *BotsService
)

func Init(cfg *config.Config) error {
	Logger = setupLogger(cfg.Env)

	var processorService = service.NewProcessor(
		Logger,
		ansmemory.NewMemoryAnswerRepository(),
		blockmemory.NewMemoryBlockRepository(),
		botmemory.NewMemoryBotRepository(),
		prtmemory.NewMemoryParticipantRepository(),
	)

	BotsAppService = NewBotsService(
		processorService,
	)

	return nil
}
