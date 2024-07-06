package domain

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type BlockRepository interface {
	Save(block *entity.Block) error
	Get(botId value.BotId, state value.State) (*entity.Block, error)
}
