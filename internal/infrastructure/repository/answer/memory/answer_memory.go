package memory

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
)

type answerMemoryRepository struct {
	m map[value.AnswerId]*entity.Answer
}

func NewMemoryAnswerRepository() interfaces.AnswerRepository {
	return &answerMemoryRepository{
		m: make(map[value.AnswerId]*entity.Answer),
	}
}

func (r *answerMemoryRepository) Save(
	_ context.Context,
	answer *entity.Answer,
) error {
	if _, ok := r.m[answer.Id]; ok {
		return interfaces.ErrAnswerAlreadyExists
	}

	r.m[answer.Id] = answer
	return nil
}

func (r *answerMemoryRepository) AnswersFrom(
	_ context.Context,
	prtId value.ParticipantId,
) ([]*entity.Answer, error) {
	answers := make([]*entity.Answer, 0)

	for _, ans := range r.m {
		if ans.Id.ParticipantId == prtId {
			answers = append(answers, ans)
		}
	}

	return answers, nil
}
