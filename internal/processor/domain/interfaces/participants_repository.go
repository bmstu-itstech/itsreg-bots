package interfaces

import (
	"context"

	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"
)

type ParticipantRepository interface {
	UpdateOrCreate(
		ctx context.Context,
		bot *bots.Bot,
		userID int64,
		updateFn func(context.Context, *bots.Participant) error,
	) error
}
