package memory

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"sync"
)

type BotMemoryRepository struct {
	m      map[value.BotId]*entity.Bot
	lastId value.BotId
	sync.RWMutex
}

func NewMemoryBotRepository() *BotMemoryRepository {
	return &BotMemoryRepository{
		m:      make(map[value.BotId]*entity.Bot),
		lastId: 0,
	}
}

func (r *BotMemoryRepository) Save(
	_ context.Context,
	bot *entity.Bot,
) (value.BotId, error) {
	r.Lock()
	defer r.Unlock()

	r.lastId++
	bot.Id = r.lastId

	r.m[bot.Id] = bot
	return r.lastId, nil
}

func (r *BotMemoryRepository) Bot(
	_ context.Context,
	id value.BotId,
) (*entity.Bot, error) {
	r.RLock()
	defer r.RUnlock()

	bot, ok := r.m[id]
	if !ok {
		return nil, interfaces.ErrBotNotFound
	}

	return bot, nil
}

func (r *BotMemoryRepository) Close() error {
	return nil
}
