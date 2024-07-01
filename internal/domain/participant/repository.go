package participant

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/entity"
)

var (
	ErrParticipantNotFound      = errors.New("participant not found")
	ErrParticipantAlreadyExists = errors.New("participant already exists")
)

type Repository interface {
	Save(prt entity.Participant) error
	Get(id entity.ParticipantId) (entity.Participant, error)
	UpdateCurrentId(id entity.ParticipantId, currentId entity.NodeId) error
}
