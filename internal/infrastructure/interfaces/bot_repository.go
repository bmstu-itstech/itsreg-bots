package interfaces

import (
	"context"
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

var (
	ErrBotNotFound = errors.New("bot not found")
)

type BotRepository interface {
	Save(context.Context, *entity.Bot) (value.BotId, error)
	Bot(context.Context, value.BotId) (*entity.Bot, error)
}
