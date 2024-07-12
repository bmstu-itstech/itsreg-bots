package memory

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
)

type participantMemoryRepository struct {
	m map[value.ParticipantId]*entity.Participant
}

func NewMemoryParticipantRepository() interfaces.ParticipantRepository {
	return &participantMemoryRepository{
		m: make(map[value.ParticipantId]*entity.Participant),
	}
}

func (r *participantMemoryRepository) Save(
	_ context.Context,
	participant *entity.Participant,
) error {
	if _, ok := r.m[participant.Id]; ok {
		return interfaces.ErrParticipantAlreadyExists
	}

	r.m[participant.Id] = participant
	return nil
}

func (r *participantMemoryRepository) Participant(
	_ context.Context,
	id value.ParticipantId,
) (*entity.Participant, error) {
	prt, ok := r.m[id]
	if !ok {
		return nil, interfaces.ErrParticipantNotFound
	}

	return prt, nil
}

func (r *participantMemoryRepository) UpdateState(
	_ context.Context,
	id value.ParticipantId,
	state value.State,
) error {
	prt, ok := r.m[id]
	if !ok {
		return interfaces.ErrParticipantNotFound
	}

	prt.Current = state

	return nil
}
