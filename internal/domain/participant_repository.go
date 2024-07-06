package domain

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

var (
	ErrParticipantNotFound = errors.New("participant not found")
)

type ParticipantRepository interface {
	Save(participant *entity.Participant) error
	Get(id value.ParticipantId) (*entity.Participant, error)
	UpdateState(id value.ParticipantId, state value.State) error
}
