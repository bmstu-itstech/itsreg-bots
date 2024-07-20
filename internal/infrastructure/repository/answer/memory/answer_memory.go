package memory

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"sync"
)

type AnswerMemoryRepository struct {
	m map[value.AnswerId]*entity.Answer
	sync.RWMutex
}

func NewMemoryAnswerRepository() *AnswerMemoryRepository {
	return &AnswerMemoryRepository{
		m: make(map[value.AnswerId]*entity.Answer),
	}
}

func (r *AnswerMemoryRepository) Save(
	_ context.Context,
	answer *entity.Answer,
) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[answer.Id]; ok {
		return interfaces.ErrAnswerExists
	}

	r.m[answer.Id] = answer
	return nil
}

func (r *AnswerMemoryRepository) AnswersFrom(
	_ context.Context,
	prtId value.ParticipantId,
) ([]*entity.Answer, error) {
	r.RLock()
	defer r.RUnlock()

	answers := make([]*entity.Answer, 0)

	for _, ans := range r.m {
		if ans.Id.ParticipantId == prtId {
			answers = append(answers, ans)
		}
	}

	return answers, nil
}

func (r *AnswerMemoryRepository) Close() error {
	return nil
}
