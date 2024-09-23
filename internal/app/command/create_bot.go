package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/app/types"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type CreateBot struct {
	AuthorUUID string

	BotUUID string
	Name    string
	Token   string

	Entries  []types.EntryPoint
	Mailings []types.Mailing
	Blocks   []types.Block
}

type CreateBotHandler decorator.CommandHandler[CreateBot]

type createBotHandler struct {
	bots bots.Repository
}

func NewCreateBotHandler(
	bots bots.Repository,

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

	mailings, err := types.MapMailingsToDomain(cmd.Mailings)
	if err != nil {
		return err
	}

	blocks, err := types.MapBlocksToDomain(cmd.Blocks)
	if err != nil {
		return err
	}

	bot, err := bots.NewBot(cmd.BotUUID, cmd.AuthorUUID, entries, mailings, blocks, cmd.Name, cmd.Token)
	if err != nil {
		return err
	}

	err = h.bots.UpdateOrCreate(ctx, bot)
	if err != nil {
		return err
	}

	return nil
}
