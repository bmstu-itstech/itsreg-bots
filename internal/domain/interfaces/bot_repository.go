package interfaces

import (
	"context"
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

var (
	ErrBotNotFound = errors.New("bot not found")
)

type BotRepository interface {
	Save(context.Context, *entity.Bot) error
	Bot(context.Context, value.BotId) (*entity.Bot, error)
}
