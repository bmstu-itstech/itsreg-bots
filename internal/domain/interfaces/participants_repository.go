package interfaces

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type ParticipantRepository interface {
	UpdateOrCreate(
		ctx context.Context,
		id int64,
		updateFn func(context.Context, *bots.Participant) (*bots.Participant, error),
	) error
}
