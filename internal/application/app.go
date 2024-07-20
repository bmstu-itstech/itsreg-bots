package application

import (
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/application/processor"
	"github.com/bmstu-itstech/itsreg-bots/internal/application/registry"
	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	pgans "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/answer/postgres"
	pgblc "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/block/postgres"
	pgbot "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/bot/postgres"
	pgprt "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/participant/postgres"
	"github.com/bmstu-itstech/itsreg-bots/pkg/endpoint"
	"log/slog"
)

type App struct {
	Logger *slog.Logger

	ansRepos interfaces.AnswerRepository
	blcRepos interfaces.BlockRepository
	botRepos interfaces.BotRepository
	prtRepos interfaces.ParticipantRepository

	Processor *processor.Processor
	Registry  *registry.Registry
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	url := endpoint.BuildPostgresConnectionString(
		endpoint.WithPostgresUser(cfg.Postgres.User),
		endpoint.WithPostgresPassword(cfg.Postgres.Pass),
		endpoint.WithPostgresHost(cfg.Postgres.Host),
		endpoint.WithPostgresPort(cfg.Postgres.Port),
		endpoint.WithPostgresDb(cfg.Postgres.DbName),
	)

	ansRepos, err := pgans.NewPostgresAnswerRepository(url)
	if err != nil {
		return nil, err
	}

	blcRepos, err := pgblc.NewPostgresBlockRepository(url)
	if err != nil {
		return nil, err
	}

	botRepos, err := pgbot.NewPostgresBotRepository(url)
	if err != nil {
		return nil, err
	}

	prtRepos, err := pgprt.NewPostgresParticipantRepository(url)
	if err != nil {
		return nil, err
	}

	proc := processor.New(log, ansRepos, blcRepos, botRepos, prtRepos)
	manager := registry.New(log, botRepos, blcRepos)

	return &App{
		Logger: log,

		ansRepos: ansRepos,
		blcRepos: blcRepos,
		botRepos: botRepos,
		prtRepos: prtRepos,

		Processor: proc,
		Registry:  manager,
	}, nil
}

func (a *App) MustClose() error {
	var ret error

	if err := a.ansRepos.Close(); err != nil {
		errors.Join(ret, err)
	}

	if err := a.blcRepos.Close(); err != nil {
		errors.Join(ret, err)
	}

	if err := a.botRepos.Close(); err != nil {
		errors.Join(ret, err)
	}

	if err := a.prtRepos.Close(); err != nil {
		errors.Join(ret, err)
	}

	return ret
}
