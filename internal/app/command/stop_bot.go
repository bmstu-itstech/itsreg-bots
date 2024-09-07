package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type StopBot struct {
	UserUUID string
	BotUUID  string
}

type StopBotHandler decorator.CommandHandler[StopBot]

type stopBotHandler struct {
	bots   bots.Repository
	runPub bots.RunnerPublisher
}

func NewStopBotHandler(
	bots bots.Repository,
	runPub bots.RunnerPublisher,

	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) StopBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

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
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	if err = bot.CanSeeBot(cmd.UserUUID); err != nil {
		return err
	}

	return h.runPub.PublishStop(ctx, cmd.BotUUID)
}
