package interfaces

import (
	"context"
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

var (
	ErrAnswerExists = errors.New("answer already exists")
)

type AnswerRepository interface {
	Close() error
	Save(context.Context, *entity.Answer) error
	AnswersFrom(context.Context, value.ParticipantId) ([]*entity.Answer, error)
}
