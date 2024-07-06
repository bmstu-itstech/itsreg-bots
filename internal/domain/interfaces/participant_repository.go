package interfaces

import (
	"context"
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

var (
	ErrParticipantAlreadyExists = errors.New("participant already exists")
	ErrParticipantNotFound      = errors.New("participant not found")
)

type ParticipantRepository interface {
	Save(context.Context, *entity.Participant) error
	Participant(context.Context, value.ParticipantId) (*entity.Participant, error)
	UpdateState(context.Context, value.ParticipantId, value.State) error
}
