package interfaces

import (
	"context"
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

var (
	ErrAnswerAlreadyExists = errors.New("answer already exists")
)

type AnswerRepository interface {
	Save(context.Context, *entity.Answer) error
	AnswersFrom(context.Context, value.ParticipantId) ([]*entity.Answer, error)
}
