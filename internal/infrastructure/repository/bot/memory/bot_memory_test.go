package memory

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/interfaces"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	"testing"
)

func TestBotMemoryRepository_Save(t *testing.T) {
	t.Run("should save the bot", func(t *testing.T) {
		ctx := context.Background()
		repos := botMemoryRepository{
			m: map[value.BotId]*entity.Bot{},
		}

		bot, err := entity.NewBot(0, "example bot", "example token", value.State(42))
		require.NoError(t, err)

		err = repos.Save(ctx, bot)
		require.NoError(t, err)

		got, ok := repos.m[bot.Id]
		require.True(t, ok, "bot should be saved")
		require.Equal(t, *bot, *got)
	})
}

func TestBotMemoryRepository_Bot(t *testing.T) {
	expected, err := entity.NewBot(
		value.BotId(42),
		"example bot",
		"example token",
		value.State(37))
	require.NoError(t, err)

	repos := botMemoryRepository{
		m: map[value.BotId]*entity.Bot{
			expected.Id: expected,
		},
	}

	t.Run("should find the bot", func(t *testing.T) {
		ctx := context.Background()

		bot, err := repos.Bot(ctx, value.BotId(42))
		require.NoError(t, err)
		require.Equal(t, *expected, *bot)
	})

	t.Run("should return error when bot does not exists", func(t *testing.T) {
		ctx := context.Background()

		_, err := repos.Bot(ctx, value.BotId(37))
		require.ErrorIs(t, err, interfaces.ErrBotNotFound)
	})
}
