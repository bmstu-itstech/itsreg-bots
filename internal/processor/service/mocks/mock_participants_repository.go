package mocks

import (
	"context"
	"sync"

	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/interfaces"
)

type participantID struct {
	BotUUID string
	UserID  int64
}

type mockParticipantsRepository struct {
	sync.Mutex
	m map[participantID]bots.Participant
}

func NewMockParticipantsRepository() interfaces.ParticipantRepository {
	return &mockParticipantsRepository{m: make(map[participantID]bots.Participant)}
}

func (r *mockParticipantsRepository) UpdateOrCreate(
	ctx context.Context,
	bot *bots.Bot,
	userID int64,
	updateFn func(ctx context.Context, prt *bots.Participant) error,
) error {
	r.Lock()
	defer r.Unlock()

	id := participantID{
		BotUUID: bot.UUID,
		UserID:  userID,
	}
	prt, ok := r.m[id]
	if !ok {
		newPrt, err := bots.NewParticipant(userID, bot.StartState)
		if err != nil {
			return err
		}
		r.m[id] = *newPrt
		prt = *newPrt
	}

	return updateFn(ctx, &prt)
}
