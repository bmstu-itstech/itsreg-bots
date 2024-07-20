package memory

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"sync"
)

type ParticipantMemoryRepository struct {
	m map[value.ParticipantId]*entity.Participant
	sync.RWMutex
}

func NewMemoryParticipantRepository() *ParticipantMemoryRepository {
	return &ParticipantMemoryRepository{
		m: make(map[value.ParticipantId]*entity.Participant),
	}
}

func (r *ParticipantMemoryRepository) Save(
	_ context.Context,
	participant *entity.Participant,
) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[participant.Id]; ok {
		return interfaces.ErrParticipantAlreadyExists
	}

	r.m[participant.Id] = participant
	return nil
}

func (r *ParticipantMemoryRepository) Participant(
	_ context.Context,
	id value.ParticipantId,
) (*entity.Participant, error) {
	r.RLock()
	defer r.RUnlock()

	prt, ok := r.m[id]
	if !ok {
		return nil, interfaces.ErrParticipantNotFound
	}

	return prt, nil
}

func (r *ParticipantMemoryRepository) UpdateState(
	_ context.Context,
	id value.ParticipantId,
	state value.State,
) error {
	r.Lock()
	defer r.Unlock()

	prt, ok := r.m[id]
	if !ok {
		return interfaces.ErrParticipantNotFound
	}

	prt.Current = state

	return nil
}

func (r *ParticipantMemoryRepository) Close() error {
	return nil
}
