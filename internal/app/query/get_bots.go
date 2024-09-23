package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type GetBots struct {
	UserUUID string
}

type GetBotsHandler decorator.QueryHandler[GetBots, []types.Bot]

type getBotsHandler struct {
	bots bots.Repository
}

func NewGetBotsHandler(
	bots bots.Repository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetBotsHandler {
	return decorator.ApplyQueryDecorators[GetBots, []types.Bot](
		getBotsHandler{bots: bots},
		logger,
		metricsClient,
	)
}

func (h getBotsHandler) Handle(ctx context.Context, query GetBots) ([]types.Bot, error) {
	bs, err := h.bots.UserBots(ctx, query.UserUUID)
	if err != nil {
		return nil, err
	}

	return types.MapBotsFromDomain(bs), nil
}
