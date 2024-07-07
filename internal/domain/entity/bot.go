package entity

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

var (
	ErrInvalidBotName       = errors.New("invalid bot name")
	ErrInvalidBotToken      = errors.New("invalid bot token")
	ErrInvalidBotStartState = errors.New("invalid bot start state")
)

type Bot struct {
	Id    value.BotId
	Name  string
	Token string
	Start value.State
}

func NewBot(id value.BotId, name string, token string, start value.State) (*Bot, error) {
	if len(name) == 0 {
		return nil, ErrInvalidBotName
	}

	if len(token) == 0 {
		return nil, ErrInvalidBotToken
	}

	if start.IsNone() {
		return nil, ErrInvalidBotStartState
	}

	return &Bot{
		Id:    id,
		Name:  name,
		Token: token,
		Start: start,
	}, nil
}
