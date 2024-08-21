package service

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/metrics"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/app/query"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/infra"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/service/mocks"
)

type CloseFunc func() error

func NewApplication() (*app.Application, CloseFunc) {
	logger := &slog.Logger{}
	metricsClient := metrics.NoOp{}

	bots := mocks.NewMockBotRepository()
	answers := mocks.NewMockAnswersRepository()
	participants := mocks.NewMockParticipantsRepository()

	actor := infra.NewTelegramActorService(
		logger, processor{h: command.NewProcessHandler(bots, participants, answers, logger, metricsClient)},
	)

	return &app.Application{
			Commands: app.Commands{
				CreateBot: command.NewCreateBotHandler(bots, logger, metricsClient),
				Process:   command.NewProcessHandler(bots, participants, answers, logger, metricsClient),
				StartBot:  command.NewStartBotHandler(bots, actor, logger, metricsClient),
				StopBot:   command.NewStopBotHandler(bots, actor, logger, metricsClient),
			},
			Queries: app.Queries{
				AllAnswers: query.NewAllAnswersHandler(bots, answers, logger, metricsClient),
			},
		}, func() error {
			return nil
		}
}

type processor struct {
	h command.ProcessHandler
}

func (p processor) Process(
	ctx context.Context, botUUID string, userID int64, text string,
) error {
	return p.h.Handle(ctx, command.Process{BotUUID: botUUID, UserId: userID, Text: text})
}
