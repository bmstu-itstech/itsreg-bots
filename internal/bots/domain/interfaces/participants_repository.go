package interfaces

import (
	"context"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
)

type ParticipantRepository interface {
	ParticipantsOfBot(ctx context.Context, botUUID string) ([]*bots.Participant, error)
	UpdateOrCreate(
		ctx context.Context,
		botUUID string,
		userID int64,
		updateFn func(context.Context, *bots.Participant) error,
	) error
}
