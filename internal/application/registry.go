package application

import (
	"github.com/zhikh23/itsreg-bots/internal/config"
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

	memoryAnswerRepos := ansmemory.NewMemoryAnswerRepository()
	memoryBlockRepos := blockmemory.NewMemoryBlockRepository()
	memoryBotRepos := botmemory.NewMemoryBotRepository()
	memoryPrtRepos := prtmemory.NewMemoryParticipantRepository()

	botsProcessorService := NewProcessor(
		Logger,
		memoryAnswerRepos,
		memoryBlockRepos,
		memoryBotRepos,
		memoryPrtRepos,
	)

	botsManagerService := NewBotsManager(
		Logger,
		memoryBotRepos,
		memoryBlockRepos,
	)

	BotsAppService = NewBotsService(
		botsProcessorService,
		botsManagerService,
	)

	return nil
}
