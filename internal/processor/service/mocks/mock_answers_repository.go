package mocks

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"
	"github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/interfaces"
	"sync"
)

type answerID struct {
	BotUUID string
	UserID  int64
	State   int
}

type mockAnswersRepository struct {
	sync.RWMutex
	m map[answerID]bots.Answer
}

func NewMockAnswersRepository() interfaces.AnswersRepository {
	return &mockAnswersRepository{m: make(map[answerID]bots.Answer)}
}

func (r *mockAnswersRepository) Upsert(
	_ context.Context,
	botUUID string,
	ans *bots.Answer,
) error {
	r.Lock()
	defer r.Unlock()

	id := answerID{
		BotUUID: botUUID,
		UserID:  ans.UserID,
		State:   ans.State,
	}
	r.m[id] = *ans

	return nil
}

func (r *mockAnswersRepository) AnswersFromParticipant(
	_ context.Context,
	botUUID string,
	userID int64,
) ([]*bots.Answer, error) {
	r.RLock()
	defer r.RUnlock()

	answers := make([]*bots.Answer, 0)

	for id, answer := range r.m {
		if id.BotUUID == botUUID && id.UserID == userID {
			answers = append(answers, &answer)
		}
	}

	return answers, nil
}
