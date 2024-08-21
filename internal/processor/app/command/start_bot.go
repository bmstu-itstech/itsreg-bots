package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/interfaces"
)

type StartBot struct {
	BotUUID string
}

type StartBotHandler decorator.CommandHandler[StartBot]

type startBotHandler struct {
	bots  interfaces.BotsRepository
	actor interfaces.ActorService
}

func NewStartBotHandler(
	bots interfaces.BotsRepository,
	actor interfaces.ActorService,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if actor == nil {
		panic("actor service is nil")
	}

	return decorator.ApplyCommandDecorators[StartBot](
		startBotHandler{bots: bots, actor: actor},
		logger,
		metricsClient,
	)
}

func (h startBotHandler) Handle(ctx context.Context, cmd StartBot) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	err = h.actor.Start(ctx, bot)
	if err != nil {
		return err
	}

	return nil
}
