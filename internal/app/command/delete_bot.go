package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type DeleteBot struct {
	AuthorUUID string
	BotUUID    string
}

type DeleteBotHandler decorator.CommandHandler[DeleteBot]

type deleteBotHandler struct {
	bots   bots.Repository
	runPub bots.RunnerPublisher
}

func NewDeleteBotHandler(
	bots bots.Repository,
	runPub bots.RunnerPublisher,
	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) DeleteBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	return decorator.ApplyCommandDecorators[DeleteBot](
		deleteBotHandler{bots, runPub},
		logger, metricsClient,
	)
}

func (h deleteBotHandler) Handle(ctx context.Context, cmd DeleteBot) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	if bot.OwnerUUID != cmd.AuthorUUID {
		return bots.ErrPermissionDenied
	}

	err = h.runPub.PublishStop(ctx, bot.UUID)
	if err != nil {
		return err
	}

	return h.bots.Delete(ctx, cmd.BotUUID)
}
