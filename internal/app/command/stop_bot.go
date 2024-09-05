package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type StopBot struct {
	BotUUID string
}

type StopBotHandler decorator.CommandHandler[StopBot]

type stopBotHandler struct {
	runPub bots.RunnerPublisher
}

func NewStopBotHandler(
	runPub bots.RunnerPublisher,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) StopBotHandler {
	if runPub == nil {
		panic("runner publisher is nil")
	}

	return decorator.ApplyCommandDecorators[StopBot](
		&stopBotHandler{runPub: runPub},
		log,
		metricsClient,
	)
}

func (h stopBotHandler) Handle(ctx context.Context, cmd StopBot) error {
	return h.runPub.PublishStop(ctx, cmd.BotUUID)
}
