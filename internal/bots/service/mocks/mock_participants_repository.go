package mocks

import (
	"context"
	"sync"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
)

type participantID struct {
	BotUUID string
	UserID  int64
}

type mockParticipantsRepository struct {
	sync.RWMutex
	m map[participantID]bots.Participant
}

func NewMockParticipantsRepository() interfaces.ParticipantRepository {
	return &mockParticipantsRepository{m: make(map[participantID]bots.Participant)}
}

func (r *mockParticipantsRepository) ParticipantsOfBot(
	ctx context.Context,
	botUUID string,
) ([]*bots.Participant, error) {
	r.RLock()
	defer r.RUnlock()

	prts := make([]*bots.Participant, 0)
	for _, p := range r.m {
		if p.BotUUID == botUUID {
			prts = append(prts, &p)
		}
	}

	return prts, nil
}

func (r *mockParticipantsRepository) UpdateOrCreate(
	ctx context.Context,
	botUUID string,
	userID int64,
	updateFn func(ctx context.Context, prt *bots.Participant) error,
) error {
	r.Lock()
	defer r.Unlock()

	id := participantID{
		BotUUID: botUUID,
		UserID:  userID,
	}
	prt, ok := r.m[id]
	if !ok {
		newPrt, err := bots.NewParticipant(botUUID, userID)
		if err != nil {
			return err
		}
		r.m[id] = *newPrt
		prt = *newPrt
	}

	err := updateFn(ctx, &prt)
	if err != nil {
		return err
	}

	r.m[id] = prt

	return nil
}
