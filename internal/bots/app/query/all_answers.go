package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
)

type AllAnswers struct {
	BotUUID string
	UserID  int64
}

type AllAnswersHandler decorator.QueryHandler[AllAnswers, Table]

type allAnswersHandler struct {
	bots    interfaces.BotsRepository
	answers interfaces.AnswersRepository
}

func NewAllAnswersHandler(
	bots interfaces.BotsRepository,
	answers interfaces.AnswersRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) AllAnswersHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if answers == nil {
		panic("answers repository is nil")
	}

	return decorator.ApplyQueryDecorators[AllAnswers, Table](
		allAnswersHandler{bots: bots, answers: answers},
		logger,
		metricsClient,
	)
}

func (h allAnswersHandler) Handle(ctx context.Context, query AllAnswers) (Table, error) {
	bot, err := h.bots.Bot(ctx, query.BotUUID)
	if err != nil {
		return Table{}, err
	}

	answers, err := h.answers.AnswersFromParticipant(ctx, query.BotUUID, query.UserID)
	if err != nil {
		return Table{}, err
	}

	table, err := bots.NewTable(bot, answers)
	if err != nil {
		return Table{}, err
	}

	return mapTableFromDomain(table), nil
}
