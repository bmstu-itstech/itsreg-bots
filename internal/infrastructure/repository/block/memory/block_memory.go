package memory

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"sync"
)

type blockId struct {
	botId value.BotId
	state value.State
}

type BlockMemoryRepository struct {
	m map[blockId]*entity.Block
	sync.RWMutex
}

func NewMemoryBlockRepository() *BlockMemoryRepository {
	return &BlockMemoryRepository{
		m: make(map[blockId]*entity.Block),
	}
}

func (r *BlockMemoryRepository) Save(
	_ context.Context,
	block *entity.Block,
) error {
	r.Lock()
	defer r.Unlock()

	id := blockId{
		botId: block.BotId,
		state: block.State,
	}

	if _, ok := r.m[id]; ok {
		return interfaces.ErrBlockAlreadyExists
	}

	r.m[id] = block
	return nil
}

func (r *BlockMemoryRepository) Block(
	_ context.Context,
	botId value.BotId,
	state value.State,
) (*entity.Block, error) {
	r.RLock()
	defer r.RUnlock()

	id := blockId{
		botId: botId,
		state: state,
	}

	block, ok := r.m[id]
	if !ok {
		return nil, interfaces.ErrBlockNotFound
	}

	return block, nil
}
