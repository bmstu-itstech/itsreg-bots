package domain

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type BotRepository interface {
	Save(bot *entity.Bot) error
	Get(id value.BotId) (*entity.Bot, error)
}
