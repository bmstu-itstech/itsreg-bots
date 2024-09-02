package service

import (
	"errors"
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/infra"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs"
	"github.com/jmoiron/sqlx"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/query"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/service/mocks"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs/handlers/slogdiscard"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/metrics"
)

type CloseFunc func() error

func NewApplication() (*app.Application, CloseFunc) {
	logger := logs.DefaultLogger()
	metricsClient := metrics.NoOp{}

	url := os.Getenv("DATABASE_URL")
	db := sqlx.MustConnect("postgres", url)

	bots := infra.NewPgBotsRepository(db)
	participants := infra.NewPgParticipantsRepository(db)

	sender, senderClose := infra.NewAmqpSenderService()

	return newApplication(logger, metricsClient, bots, participants, sender), func() error {
		var err error
		err = errors.Join(err, db.Close())
		err = errors.Join(err, senderClose())
		return err
	}
}

func NewComponentTestApplication() (*app.Application, message.Subscriber) {
	logger := slogdiscard.NewDiscardLogger()
	metricsClient := metrics.NoOp{}

	bots := mocks.NewMockBotRepository()
	participants := mocks.NewMockParticipantsRepository()

	sender, subscriber := mocks.NewMockSenderService()

	return newApplication(logger, metricsClient, bots, participants, sender), subscriber
}

func newApplication(
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
	bots interfaces.BotsRepository,
	participants interfaces.ParticipantRepository,
	sender interfaces.SenderService,
) *app.Application {
	return &app.Application{
		Commands: app.Commands{
			CreateBot: command.NewCreateBotHandler(bots, logger, metricsClient),
			Entry:     command.NewEntryHandler(bots, participants, sender, logger, metricsClient),
			Process:   command.NewProcessHandler(bots, participants, sender, logger, metricsClient),
		},
		Queries: app.Queries{
			AllAnswers: query.NewGetAnswersTableHandler(bots, participants, logger, metricsClient),
			GetBot:     query.NewGetBotHandler(bots, logger, metricsClient),
		},
	}
}
