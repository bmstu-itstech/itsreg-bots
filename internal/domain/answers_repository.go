package domain

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type AnswersRepository interface {
	Save(*entity.Answer) error
	Get(id value.AnswerId) (*entity.Answer, error)
}
