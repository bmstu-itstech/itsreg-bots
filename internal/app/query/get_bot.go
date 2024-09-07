package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type GetBot struct {
	UserUUID string
	BotUUID  string
}

type GetBotHandler decorator.QueryHandler[GetBot, types.Bot]

type getBotHandler struct {
	bots bots.Repository
}

func NewGetBotHandler(
	bots bots.Repository,

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

	if err = bot.CanSeeBot(query.UserUUID); err != nil {
		return types.Bot{}, err
	}

	return types.MapBotFromDomain(bot), nil
}
