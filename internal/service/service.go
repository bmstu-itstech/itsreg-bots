package service

import (
	"errors"
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"

	"github.com/bmstu-itstech/itsreg-bots/internal/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/app/query"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs/handlers/slogdiscard"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/metrics"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/infra"
	"github.com/bmstu-itstech/itsreg-bots/internal/service/mocks"
)

func NewApplication() (
	app *app.Application,
	msgCh <-chan *message.Message,
	runCh <-chan *message.Message,
	close func() error,
) {
	logger := logs.DefaultLogger()
	metricsClient := metrics.NoOp{}

	url := os.Getenv("DATABASE_URI")
	db := sqlx.MustConnect("postgres", url)

	botsR := infra.NewPgBotsRepository(db)
	participants := infra.NewPgParticipantsRepository(db)

	msgPub, msgCh, senderClose := infra.NewNATSMessagesPublisher()
	runPub, runCh, senderClose := infra.NewNATSRunnerPublisher()

	return newApplication(
			logger, metricsClient, botsR, participants, msgPub, runPub,
		), msgCh, runCh, func() error {
			var err error
			err = errors.Join(err, db.Close())
			err = errors.Join(err, senderClose())
			return err
		}
}

func NewComponentTestApplication() (
	app *app.Application,
	msgCh <-chan *message.Message,
	runCh <-chan *message.Message,
) {
	logger := slogdiscard.NewDiscardLogger()
	metricsClient := metrics.NoOp{}

	botsR := mocks.NewMockBotRepository()
	participants := mocks.NewMockParticipantsRepository()

	msgPub, msgCh := mocks.NewMockMessagesPublisher()
	runPub, runCh := mocks.NewMockRunnerPublisher()

	return newApplication(
		logger, metricsClient, botsR, participants, msgPub, runPub,
	), msgCh, runCh
}

func newApplication(
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
	bots bots.Repository,
	participants bots.ParticipantRepository,
	msgPub bots.MessagesPublisher,
	runPub bots.RunnerPublisher,
) *app.Application {
	return &app.Application{
		Commands: app.Commands{
			CreateBot:    command.NewCreateBotHandler(bots, logger, metricsClient),
			StartBot:     command.NewStartBotHandler(bots, runPub, logger, metricsClient),
			StopBot:      command.NewStopBotHandler(bots, runPub, logger, metricsClient),
			UpdateStatus: command.NewUpdateStatusHandler(bots, logger, metricsClient),
			Entry:        command.NewEntryHandler(bots, participants, msgPub, logger, metricsClient),
			Process:      command.NewProcessHandler(bots, participants, msgPub, logger, metricsClient),
			StartMailing: command.NewStartMailingHandler(bots, participants, msgPub, logger, metricsClient),
		},
		Queries: app.Queries{
			AllAnswers: query.NewGetAnswersTableHandler(bots, participants, logger, metricsClient),
			GetBot:     query.NewGetBotHandler(bots, logger, metricsClient),
			GetBots:    query.NewGetBotsHandler(bots, logger, metricsClient),
		},
	}
}
