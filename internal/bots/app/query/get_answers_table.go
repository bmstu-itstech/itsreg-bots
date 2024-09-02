package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
)

type GetAnswersTable struct {
	BotUUID string
}

type GetAnswersTableHandler decorator.QueryHandler[GetAnswersTable, types.AnswersTable]

type answersHandler struct {
	bots         interfaces.BotsRepository
	participants interfaces.ParticipantRepository
}

func NewGetAnswersTableHandler(
	bots interfaces.BotsRepository,
	participants interfaces.ParticipantRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetAnswersTableHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if participants == nil {
		panic("participants repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetAnswersTable, types.AnswersTable](
		answersHandler{bots: bots, participants: participants},
		logger,
		metricsClient,
	)
}

func (h answersHandler) Handle(ctx context.Context, query GetAnswersTable) (types.AnswersTable, error) {
	bot, err := h.bots.Bot(ctx, query.BotUUID)
	if err != nil {
		return types.AnswersTable{}, err
	}

	prts, err := h.participants.ParticipantsOfBot(ctx, query.BotUUID)
	if err != nil {
		return types.AnswersTable{}, err
	}

	table := bots.NewAnswersTable(bot, prts)

	return types.MapAnswersTableFromDomain(table), nil
}
