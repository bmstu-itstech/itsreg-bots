package interfaces

import (
	"context"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
)

type AnswersRepository interface {
	Upsert(ctx context.Context, botUUID string, answer *bots.Answer) error
	AnswersFromParticipant(ctx context.Context, botUUID string, userID int64) ([]*bots.Answer, error)
}
