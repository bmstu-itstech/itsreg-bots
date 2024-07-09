package memory

import (
	"context"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	"github.com/zhikh23/itsreg-bots/internal/infrastructure/interfaces"
)

type blockId struct {
	botId value.BotId
	state value.State
}

type blockMemoryRepository struct {
	m map[blockId]*entity.Block
}

func NewMemoryBlockRepository() interfaces.BlockRepository {
	return &blockMemoryRepository{
		m: make(map[blockId]*entity.Block),
	}
}

func (r *blockMemoryRepository) Save(
	_ context.Context,
	block *entity.Block,
) error {
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

func (r *blockMemoryRepository) Block(
	_ context.Context,
	botId value.BotId,
	state value.State,
) (*entity.Block, error) {
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
