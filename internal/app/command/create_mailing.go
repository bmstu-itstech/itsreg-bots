package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type CreateMailing struct {
	AuthorUUID string
	BotUUID    string

	MailingName   string
	RequiredState int

	EntryPoint types.EntryPoint
	Blocks     []types.Block
}

type CreateMailingHandler decorator.CommandHandler[CreateMailing]

type createMailingHandler struct {
	bots bots.Repository
}

func NewCreateMailingHandler(
	bots bots.Repository,
	log *slog.Logger,
	metricsClient decorator.MetricsClient,
) CreateMailingHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	return decorator.ApplyCommandDecorators[CreateMailing](
		&createMailingHandler{bots},
		log, metricsClient,
	)
}

func (h createMailingHandler) Handle(ctx context.Context, cmd CreateMailing) error {
	return h.bots.Update(ctx, cmd.BotUUID, func(innerCtx context.Context, bot *bots.Bot) error {
		var err error
		if err = bot.CanSeeBot(cmd.AuthorUUID); err != nil {
			return err
		}

		entry, err := types.MapEntryPointToDomain(cmd.EntryPoint)
		if err != nil {
			return err
		}

		blocks, err := types.MapBlocksToDomain(cmd.Blocks)
		if err != nil {
			return err
		}

		return bot.AddMailing(cmd.MailingName, cmd.RequiredState, entry, blocks)
	})
}
