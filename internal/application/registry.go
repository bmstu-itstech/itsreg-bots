package application

import (
	botmemory "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/bot/memory"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	ansmemory "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/answer/memory"
	blockmemory "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/block/memory"
	_ "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/bot/memory"
	prtmemory "github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/repository/participant/memory"
)

var (
	Logger         *slog.Logger
	BotsAppService *BotsService
)

func Init(cfg *config.Config) error {
	Logger = setupLogger(cfg.Env)

	answerRepos := ansmemory.NewMemoryAnswerRepository()
	blockRepos := blockmemory.NewMemoryBlockRepository()
	prtRepos := prtmemory.NewMemoryParticipantRepository()
	botRepos := botmemory.NewMemoryBotRepository()

	//postgresBotRepos, err := botpostgres.NewPostgresBotRepository(endpoint.PostgresConnectionString(
	//	cfg.Postgres.Host, fmt.Sprint(cfg.Postgres.Port), cfg.Postgres.User, cfg.Postgres.Password,
	//	cfg.Postgres.DbName, endpoint.SslModeDisabled,
	//))
	//if err != nil {
	//	return err
	//}

	botsProcessorService := NewProcessor(
		Logger,
		answerRepos,
		blockRepos,
		botRepos,
		prtRepos,
	)

	botsManagerService := NewBotsManager(
		Logger,
		botRepos,
		blockRepos,
	)

	BotsAppService = NewBotsService(
		botsProcessorService,
		botsManagerService,
	)

	return nil
}
