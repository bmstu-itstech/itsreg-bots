package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type UpdateStatus struct {
	BotUUID string
	Status  string
}

type UpdateStatusHandler decorator.CommandHandler[UpdateStatus]

type updateStatusHandler struct {
	bots bots.Repository
}

func NewUpdateStatusHandler(
	bots bots.Repository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) UpdateStatusHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	return decorator.ApplyCommandDecorators[UpdateStatus](
		updateStatusHandler{bots: bots},
		logger,
		metricsClient,
	)
}

func (h updateStatusHandler) Handle(ctx context.Context, cmd UpdateStatus) error {
	st, err := bots.NewStatusFromString(cmd.Status)
	if err != nil {
		return err
	}

	return h.bots.UpdateStatus(ctx, cmd.BotUUID, st)
}
