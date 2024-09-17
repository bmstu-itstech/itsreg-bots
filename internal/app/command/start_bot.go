package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type StartBot struct {
	AuthorUUID string
	BotUUID    string
}

type StartBotHandler decorator.CommandHandler[StartBot]

type startBotHandler struct {
	bots   bots.Repository
	runPub bots.RunnerPublisher
}

func NewStartBotHandler(
	bots bots.Repository,
	runPub bots.RunnerPublisher,

	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if runPub == nil {
		panic("runner publisher is nil")
	}

	return decorator.ApplyCommandDecorators[StartBot](
		&startBotHandler{bots: bots, runPub: runPub},
		log,
		metricsClient,
	)
}

func (h startBotHandler) Handle(ctx context.Context, cmd StartBot) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	if err = bot.CanSeeBot(cmd.AuthorUUID); err != nil {
		return err
	}

	return h.runPub.PublishStart(ctx, cmd.BotUUID)
}
