package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type StartMailing struct {
	AuthorUUID string
	BotUUID    string
	EntryKey   string
}

type StartMailingHandler decorator.CommandHandler[StartMailing]

type startMailingHandler struct {
	bots         bots.Repository
	participants bots.ParticipantRepository
	msgPublisher bots.MessagesPublisher
}

func NewStartMailingHandler(
	bots bots.Repository,
	participants bots.ParticipantRepository,
	msgPublisher bots.MessagesPublisher,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) StartMailingHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if participants == nil {
		panic("participants repository is nil")
	}

	if msgPublisher == nil {
		panic("message publisher is nil")
	}

	return decorator.ApplyCommandDecorators[StartMailing](
		&startMailingHandler{bots, participants, msgPublisher},
		logger,
		metricsClient,
	)
}

func (h *startMailingHandler) Handle(ctx context.Context, cmd StartMailing) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	mailing, err := bot.Mailing(cmd.EntryKey)
	if err != nil {
		return err
	}

	prts, err := h.participants.ParticipantsOfBot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	prts = filterParticipants(prts, func(prt *bots.Participant) bool {
		return prt.HasAnswer(mailing.RequireState)
	})

	for _, prt := range prts {
		err := h.participants.UpdateOrCreate(
			ctx,
			cmd.BotUUID, prt.UserID,
			func(innerCtx context.Context, prt *bots.Participant,
			) error {
				messages, err := bot.Entry(prt, cmd.EntryKey)
				if err != nil {
					return err
				}

				for _, message := range messages {
					err = h.msgPublisher.Publish(ctx, cmd.BotUUID, prt.UserID, message)
					if err != nil {
						return err
					}
				}

				return nil
			})
		if err != nil {
			return err
		}
	}

	return nil
}

func filterParticipants(prts []*bots.Participant, predicate func(prt *bots.Participant) bool) []*bots.Participant {
	result := make([]*bots.Participant, 0, len(prts))
	for _, p := range prts {
		if predicate(p) {
			result = append(result, p)
		}
	}
	return result
}
