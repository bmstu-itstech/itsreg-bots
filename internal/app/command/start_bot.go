package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type StartBot struct {
	BotUUID string
}

type StartBotHandler decorator.CommandHandler[StartBot]

type startBotHandler struct {
	runPub bots.RunnerPublisher
}

func NewStartBotHandler(
	runPub bots.RunnerPublisher,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartBotHandler {
	if runPub == nil {
		panic("runner publisher is nil")
	}

	return decorator.ApplyCommandDecorators[StartBot](
		&startBotHandler{runPub: runPub},
		log,
		metricsClient,
	)
}

func (h startBotHandler) Handle(ctx context.Context, cmd StartBot) error {
	return h.runPub.PublishStart(ctx, cmd.BotUUID)
}
