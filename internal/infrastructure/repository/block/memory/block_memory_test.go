package memory

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlockMemoryRepository_Save(t *testing.T) {
	t.Parallel()

	t.Run("should save a block", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repos := BlockMemoryRepository{
			m: map[blockId]*entity.Block{},
		}

		botId := value.BotId(37)
		node, err := value.NewMessageNode(value.State(1), value.State(2))
		block, err := entity.NewBlock(node, botId, "message block", "some text")
		require.NoError(t, err)

		err = repos.Save(ctx, block)
		require.NoError(t, err)

		got, ok := repos.m[blockId{
			botId: block.BotId,
			state: block.State,
		}]
		require.True(t, ok)
		require.Equal(t, *block, *got)
	})

	t.Run("should return error when answer already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repos := BlockMemoryRepository{
			m: map[blockId]*entity.Block{},
		}

		const botId = value.BotId(37)

		node, err := value.NewMessageNode(value.State(1), value.State(2))
		block, err := entity.NewBlock(node, botId, "message block", "some text")
		require.NoError(t, err)

		err = repos.Save(ctx, block)
		require.NoError(t, err)

		node, err = value.NewQuestionNode(value.State(1), value.State(2))
		block, err = entity.NewBlock(node, botId, "question block", "some question")
		require.NoError(t, err)

		err = repos.Save(ctx, block)
		require.ErrorIs(t, err, interfaces.ErrBlockAlreadyExists)
	})
}

func TestBlockMemoryRepository_Block(t *testing.T) {
	t.Parallel()

	var blocks []*entity.Block

	const botId = value.BotId(37)

	node, err := value.NewMessageNode(value.State(1), value.State(2))
	block, err := entity.NewBlock(node, value.BotId(37), "message block", "some text")
	require.NoError(t, err)
	blocks = append(blocks, block)

	node, err = value.NewQuestionNode(value.State(2), value.State(3))
	block, err = entity.NewBlock(node, botId, "question block", "some question")
	require.NoError(t, err)
	blocks = append(blocks, block)

	repos := BlockMemoryRepository{
		m: map[blockId]*entity.Block{},
	}
	for _, block := range blocks {
		id := blockId{
			botId: block.BotId,
			state: block.State,
		}
		repos.m[id] = block
	}

	t.Run("should find block", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		got, err := repos.Block(ctx, botId, value.State(2))
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, *blocks[1], *got)
	})

	t.Run("should return error when block does not exists", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		_, err := repos.Block(ctx, botId, value.State(3))
		require.ErrorIs(t, err, interfaces.ErrBlockNotFound)
	})
}
