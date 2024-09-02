package command

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/app/types"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
)

type CreateBot struct {
	BotUUID string
	Name    string
	Token   string

	Entries []types.EntryPoint
	Blocks  []types.Block
}

type CreateBotHandler decorator.CommandHandler[CreateBot]

type createBotHandler struct {
	bots interfaces.BotsRepository
}

func NewCreateBotHandler(
	bots interfaces.BotsRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) CreateBotHandler {
	if bots == nil {
		panic("bots repository is nil")
	}

	return decorator.ApplyCommandDecorators[CreateBot](
		createBotHandler{bots: bots},
		logger,
		metricsClient,
	)
}

func (h createBotHandler) Handle(ctx context.Context, cmd CreateBot) error {
	entries, err := types.MapEntriesToDomain(cmd.Entries)
	if err != nil {
		return err
	}

	blocks, err := types.MapBlocksToDomain(cmd.Blocks)
	if err != nil {
		return err
	}

	bot, err := bots.NewBot(cmd.BotUUID, entries, blocks, cmd.Name, cmd.Token)
	if err != nil {
		return err
	}

	err = h.bots.Save(ctx, bot)
	if err != nil {
		return err
	}

	return nil
}
