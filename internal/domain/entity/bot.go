package entity

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

var (
	ErrInvalidBot = errors.New("invalid bot")
)

type Bot struct {
	Id    value.BotId
	Name  string
	Token string
	Start value.State
}

func NewBot(id value.BotId, name string, token string, start value.State) (*Bot, error) {
	if len(name) == 0 {
		return nil, ErrInvalidBot
	}

	if len(token) == 0 {
		return nil, ErrInvalidBot
	}

	if start.IsNone() {
		return nil, ErrInvalidBot
	}

	return &Bot{
		Id:    id,
		Name:  name,
		Token: token,
		Start: start,
	}, nil
}
