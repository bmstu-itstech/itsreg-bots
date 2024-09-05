package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type Process struct {
	BotUUID string
	UserID  int64
	Text    string
}

type ProcessHandler decorator.CommandHandler[Process]

type processHandler struct {
	bots         bots.Repository
	participants bots.ParticipantRepository
	msgPublisher bots.MessagesPublisher
}

func NewProcessHandler(
	bots bots.Repository,
	participants bots.ParticipantRepository,
	msgPublisher bots.MessagesPublisher,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) ProcessHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if participants == nil {
		panic("participants repository is nil")
	}

	if msgPublisher == nil {
		panic("message publisher is nil")
	}

	return decorator.ApplyCommandDecorators[Process](
		processHandler{bots: bots, participants: participants, msgPublisher: msgPublisher},
		logger,
		metricsClient,
	)
}

func (h processHandler) Handle(ctx context.Context, cmd Process) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	return h.participants.UpdateOrCreate(ctx, cmd.BotUUID, cmd.UserID, func(
		innerCtx context.Context, prt *bots.Participant,
	) error {
		messages, err := bot.Process(prt, cmd.Text)
		if err != nil {
			return err
		}

		for _, message := range messages {
			err = h.msgPublisher.Publish(innerCtx, cmd.BotUUID, cmd.UserID, message)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
