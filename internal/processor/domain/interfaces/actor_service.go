package interfaces

import (
	"context"

	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"
)

type ActorService interface {
	Start(ctx context.Context, bot *bots.Bot) error
	Stop(ctx context.Context, bot *bots.Bot) error
}
