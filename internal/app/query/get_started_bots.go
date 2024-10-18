package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type GetStartedBots struct{}

type GetStartedBotsHandler decorator.QueryHandler[GetStartedBots, []types.Bot]

type getStartedBotsHandler struct {
	bots bots.Repository
}

func NewGetStartedBotsHandler(
	bots bots.Repository,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetStartedBotsHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	return decorator.ApplyQueryDecorators(
		getStartedBotsHandler{bots},
		logger, metricsClient,
	)
}

func (h getStartedBotsHandler) Handle(ctx context.Context, _ GetStartedBots) ([]types.Bot, error) {
	bs, err := h.bots.BotsWithStatus(ctx, bots.Started)
	if err != nil {
		return nil, err
	}

	return types.MapBotsFromDomain(bs), nil
}
