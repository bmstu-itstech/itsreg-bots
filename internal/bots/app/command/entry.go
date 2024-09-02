package command

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"log/slog"
)

type Entry struct {
	BotUUID string
	UserID  int64
	Key     string
}

type EntryHandler decorator.CommandHandler[Entry]

type entryHandler struct {
	bots         interfaces.BotsRepository
	participants interfaces.ParticipantRepository
	sender       interfaces.SenderService
}

func NewEntryHandler(
	bots interfaces.BotsRepository,
	participants interfaces.ParticipantRepository,
	sender interfaces.SenderService,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) EntryHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if participants == nil {
		panic("participants repository is nil")
	}

	if sender == nil {
		panic("sender service is nil")
	}

	return decorator.ApplyCommandDecorators[Entry](
		entryHandler{bots: bots, participants: participants, sender: sender},
		logger,
		metricsClient,
	)
}

func (h entryHandler) Handle(ctx context.Context, cmd Entry) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	return h.participants.UpdateOrCreate(ctx, cmd.BotUUID, cmd.UserID, func(
		innerCtx context.Context, prt *bots.Participant,
	) error {
		messages, err := bot.Entry(prt, cmd.Key)
		if err != nil {
			return err
		}

		for _, message := range messages {
			err = h.sender.Send(innerCtx, cmd.BotUUID, cmd.UserID, message)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
