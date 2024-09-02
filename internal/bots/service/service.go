package service

import (
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/command"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/query"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/service/mocks"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/logs/handlers/slogdiscard"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/metrics"
)

func NewComponentTestApplication() (*app.Application, message.Subscriber) {
	logger := slogdiscard.NewDiscardLogger()
	metricsClient := metrics.NoOp{}

	bots := mocks.NewMockBotRepository()
	participants := mocks.NewMockParticipantsRepository()

	sender, subscriber := mocks.NewMockSenderService()

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
	}, subscriber
}
