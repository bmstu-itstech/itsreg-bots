package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type GetAnswersTable struct {
	BotUUID string
}

type GetAnswersTableHandler decorator.QueryHandler[GetAnswersTable, types.AnswersTable]

type answersHandler struct {
	bots         bots.Repository
	participants bots.ParticipantRepository
}

func NewGetAnswersTableHandler(
	bots bots.Repository,
	participants bots.ParticipantRepository,

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
