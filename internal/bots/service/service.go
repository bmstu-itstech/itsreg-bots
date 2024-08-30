package service

import (
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/query"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/service/mocks"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/metrics"
)

type CloseFunc func() error

func NewApplication() (*app.Application, CloseFunc) {
	logger := &slog.Logger{}
	metricsClient := metrics.NoOp{}

	bots := mocks.NewMockBotRepository()
	participants := mocks.NewMockParticipantsRepository()
	answers := mocks.NewMockAnswersRepository()
	sender := mocks.NewMockSenderService()

	return &app.Application{
			Commands: app.Commands{
				CreateBot: command.NewCreateBotHandler(bots, logger, metricsClient),
				Process:   command.NewProcessHandler(bots, participants, answers, sender, logger, metricsClient),
			},
			Queries: app.Queries{
				AllAnswers: query.NewAllAnswersHandler(bots, answers, logger, metricsClient),
				GetBot:     query.NewGetBotHandler(bots, logger, metricsClient),
			},
		}, func() error {
			return nil
		}
}
