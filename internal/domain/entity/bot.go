package entity

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type Bot struct {
	Id    value.BotId
	Name  string
	Token string
	Start value.State
}
