package bots

import (
	"context"
)

type ParticipantRepository interface {
	ParticipantsOfBot(ctx context.Context, botUUID string) ([]*Participant, error)
	UpdateOrCreate(
		ctx context.Context,
		botUUID string,
		userID int64,
		updateFn func(context.Context, *Participant) error,
	) error
}
