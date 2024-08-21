package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/interfaces"
)

type StopBot struct {
	BotUUID string
}

type StopBotHandler decorator.CommandHandler[StopBot]

type stopBotHandler struct {
	bots interfaces.BotsRepository
	tg   interfaces.ActorService
}

func NewStopBotHandler(
	bots interfaces.BotsRepository,
	tg interfaces.ActorService,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) StopBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if tg == nil {
		panic("tg service is nil")
	}

	return decorator.ApplyCommandDecorators[StopBot](
		stopBotHandler{bots: bots, tg: tg},
		logger,
		metricsClient,
	)
}

func (h stopBotHandler) Handle(ctx context.Context, cmd StopBot) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	err = h.tg.Start(ctx, bot)
	if err != nil {
		return err
	}

	return nil
}
