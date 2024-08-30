package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
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
	sender       interfaces.SenderService
}

func NewProcessHandler(
	bots interfaces.BotsRepository,
	participants interfaces.ParticipantRepository,
	answers interfaces.AnswersRepository,
	sender interfaces.SenderService,

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

	if sender == nil {
		panic("sender service is nil")
	}

	return decorator.ApplyCommandDecorators[Process](
		processHandler{bots: bots, participants: participants, answers: answers, sender: sender},
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
		innerCtx context.Context, prt *bots.Participant,
	) error {
		blocks, ans, err := bot.Process(prt, cmd.Text)
		if err != nil {
			return err
		}

		if ans != nil {
			err = h.answers.Upsert(innerCtx, bot.UUID, ans)
			if err != nil {
				return err
			}
		}

		for _, block := range blocks {
			err = h.sender.Send(innerCtx, &block)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
