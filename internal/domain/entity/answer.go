package entity

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

var (
	ErrInvalidAnswer = errors.New("invalid answer")
)

type Answer struct {
	Id    value.AnswerId
	Value string
}

func NewAnswer(id value.AnswerId, value string) (*Answer, error) {
	if len(value) == 0 {
		return nil, ErrInvalidAnswer
	}

	return &Answer{
		Id:    id,
		Value: value,
	}, nil
}
