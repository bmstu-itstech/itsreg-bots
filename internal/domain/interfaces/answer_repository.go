package interfaces

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type AnswerRepository interface {
	Save(ctx context.Context, answer *bots.Answer) error
	AnswersFromParticipant(ctx context.Context, prtID int64) ([]*bots.Answer, error)
}
