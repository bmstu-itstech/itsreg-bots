package memory

import (
	"context"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/interfaces"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type botMemoryRepository struct {
	m      map[value.BotId]*entity.Bot
	lastId value.BotId
}

func NewMemoryBotRepository() interfaces.BotRepository {
	return &botMemoryRepository{
		m:      make(map[value.BotId]*entity.Bot),
		lastId: 0,
	}
}

func (r *botMemoryRepository) Save(
	_ context.Context,
	bot *entity.Bot,
) error {
	r.lastId++
	bot.Id = r.lastId

	r.m[bot.Id] = bot
	return nil
}

func (r *botMemoryRepository) Bot(
	_ context.Context,
	id value.BotId,
) (*entity.Bot, error) {
	bot, ok := r.m[id]
	if !ok {
		return nil, interfaces.ErrBotNotFound
	}

	return bot, nil
}
