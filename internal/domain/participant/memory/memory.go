package memory

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/participant"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/objects"
	"sync"
)

type Repository struct {
	participants map[entity.ParticipantId]entity.Participant
	sync.Mutex
}

func New() *Repository {
	return &Repository{
		participants: make(map[entity.ParticipantId]entity.Participant),
	}
}

func (r *Repository) Save(prt entity.Participant) error {
	if _, ok := r.participants[prt.ParticipantId]; ok {
		return participant.ErrParticipantAlreadyExists
	}

	r.Lock()
	r.participants[prt.ParticipantId] = prt
	r.Unlock()

	return nil
}

func (r *Repository) Get(id entity.ParticipantId) (entity.Participant, error) {
	if prt, ok := r.participants[id]; ok {
		return prt, nil
	}

	return entity.Participant{}, participant.ErrParticipantNotFound
}

func (r *Repository) UpdateCurrentId(id entity.ParticipantId, currentId objects.State) error {
	prt, ok := r.participants[id]

	if !ok {
		return participant.ErrParticipantNotFound
	}

	prt.State = currentId

	r.Lock()
	r.participants[id] = prt
	r.Unlock()

	return nil
}
