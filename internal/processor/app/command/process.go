package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/interfaces"
)

type Process struct {
	BotUUID string

	UserId int64
	Text   string
}

type ProcessHandler decorator.CommandHandler[Process]

type processHandler struct {
	bots         interfaces.BotsRepository
	participants interfaces.ParticipantRepository
	answers      interfaces.AnswersRepository
}

func NewProcessHandler(
	bots interfaces.BotsRepository,
	participants interfaces.ParticipantRepository,
	answers interfaces.AnswersRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) ProcessHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	if participants == nil {
		panic("participants repository is nil")
	}

	if answers == nil {
		panic("answers repository is nil")
	}

	return decorator.ApplyCommandDecorators[Process](
		processHandler{bots: bots, participants: participants, answers: answers},
		logger,
		metricsClient,
	)
}

func (h processHandler) Handle(ctx context.Context, cmd Process) error {
	bot, err := h.bots.Bot(ctx, cmd.BotUUID)
	if err != nil {
		return err
	}

	return h.participants.UpdateOrCreate(ctx, bot, cmd.UserId, func(
		ctx context.Context, prt *bots.Participant,
	) error {
		// TODO: public messages from bot in query
		_, ans, err := bot.Process(prt, cmd.Text)
		if err != nil {
			return err
		}

		if ans != nil {
			err = h.answers.Upsert(ctx, bot.UUID, ans)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
