package query

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"log/slog"
)

type GetBot struct {
	BotUUID string
}

type GetBotHandler decorator.QueryHandler[GetBot, types.Bot]

type getBotHandler struct {
	bots interfaces.BotsRepository
}

func NewGetBotHandler(
	bots interfaces.BotsRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetBotHandler {
	return decorator.ApplyQueryDecorators[GetBot, types.Bot](
		getBotHandler{bots: bots},
		logger,
		metricsClient,
	)
}

func (h getBotHandler) Handle(ctx context.Context, query GetBot) (types.Bot, error) {
	bot, err := h.bots.Bot(ctx, query.BotUUID)
	if err != nil {
		return types.Bot{}, err
	}

	return types.MapBotFromDomain(bot), nil
}
