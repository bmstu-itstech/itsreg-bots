package application

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/application/processor"
	"github.com/bmstu-itstech/itsreg-bots/internal/application/registry"
	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	memans "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/answer/memory"
	memblc "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/block/memory"
	membot "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/bot/memory"
	memprt "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/participant/memory"
	"log/slog"
)

type App struct {
	Logger    *slog.Logger
	Processor *processor.Processor
	Registry  *registry.Registry
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	ansRepos := memans.NewMemoryAnswerRepository()
	blcRepos := memblc.NewMemoryBlockRepository()
	botRepos := membot.NewMemoryBotRepository()
	prtRepos := memprt.NewMemoryParticipantRepository()

	proc := processor.New(log, ansRepos, blcRepos, botRepos, prtRepos)

	manager := registry.New(log, botRepos, blcRepos)

	return &App{
		Logger:    log,
		Processor: proc,
		Registry:  manager,
	}, nil
}
