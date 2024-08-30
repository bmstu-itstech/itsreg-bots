package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/internal/common/decorator"
)

type CreateBot struct {
	BotUUID string

	Blocks     []Block
	StartState int

	Name  string
	Token string
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
	blocks, err := mapBlocksToDomain(cmd.Blocks)
	if err != nil {
		return err
	}

	bot, err := bots.NewBot(cmd.BotUUID, blocks, cmd.StartState, cmd.Name, cmd.Token)
	if err != nil {
		return err
	}

	err = h.bots.Save(ctx, bot)
	if err != nil {
		return err
	}

	return nil
}
