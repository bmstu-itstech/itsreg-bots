package application_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/application"
	"github.com/zhikh23/itsreg-bots/internal/application/dto"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	blockmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/block/memory"
	botmemory "github.com/zhikh23/itsreg-bots/internal/infrastructure/repository/bot/memory"
	"github.com/zhikh23/itsreg-bots/pkg/logger/handlers/slogdiscard"
	"testing"
)

func TestBotsManager_Create(t *testing.T) {
	type args struct {
		name   string
		token  string
		start  uint64
		blocks []dto.Block
	}

	t.Run("should create a bot", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example name",
			token: "example token",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    2,
					Default: 2,
					Options: []dto.Option{},
					Title:   "Question 1",
					Text:    "Block 1",
				},
				{
					State:   2,
					Type:    1,
					Default: 0,
					Options: []dto.Option{},
					Title:   "Message 2",
					Text:    "Block 2",
				},
			},
		}

		id, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.NoError(t, err)
		require.NotZero(t, id)
	})

	t.Run("should return an error if bot name is empty", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "", // <----- should be not empty
			token: "example token",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    2,
					Default: 0,
					Options: []dto.Option{},
					Title:   "Question 1",
					Text:    "Block 1",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, entity.ErrInvalidBot)
	})

	t.Run("should return an error if bot token is empty", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "", // <------ should be not empty
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    2,
					Default: 0,
					Options: []dto.Option{},
					Title:   "Question 1",
					Text:    "Block 1",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, entity.ErrInvalidBot)
	})

	t.Run("should return an error if bot has a non-existent block", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    2,
					Default: 2, // <------ non-existent block's state
					Options: []dto.Option{},
					Title:   "Question 1",
					Text:    "Block 1",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, application.ErrNodeNotFound)
	})

	t.Run("should return an error if bot has a cycle", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    3,
					Default: 0,
					Options: []dto.Option{
						{
							Text: "To 3",
							Next: 3,
						},
						{
							Text: "To 2",
							Next: 2, // <---- cycle here
						},
					},
					Title: "Selection 1",
					Text:  "Block 1",
				},
				{
					State:   1,
					Type:    3,
					Default: 0,
					Options: []dto.Option{
						{
							Text: "To 1",
							Next: 1, // <---- cycle here
						},
						{
							Text: "To 3",
							Next: 3,
						},
					},
					Title: "Selection 2",
					Text:  "Block 2",
				},
				{
					State:   1,
					Type:    2,
					Default: 1,
					Options: []dto.Option{},
					Title:   "Question 3",
					Text:    "Block 3",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, application.ErrRecursionFound)
	})

	t.Run("should return an error if bot has an unused block", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    2,
					Default: 0,
					Options: []dto.Option{},
					Title:   "Question 1",
					Text:    "Block 1",
				},
				{
					State:   2, // <----- does not have links here
					Type:    2,
					Default: 0,
					Options: []dto.Option{},
					Title:   "Question 2",
					Text:    "Block 2",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, application.ErrUnusedBlockFound)
	})

	t.Run("should return an error if bot has invalid message block", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "example token",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    1,
					Default: 2,
					Options: []dto.Option{ // <----- should be empty
						{
							Text: "Option???",
							Next: 2,
						},
					},
					Title: "Message 1",
					Text:  "Block 1",
				},
				{
					State:   2,
					Type:    2,
					Default: 0,
					Options: []dto.Option{},
					Title:   "Question 2",
					Text:    "Block 2",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, value.ErrInvalidNode)
	})

	t.Run("should return an error if bot has invalid question block", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    2,
					Default: 2,
					Options: []dto.Option{ // <----- should be empty
						{
							Text: "Option???",
							Next: 2,
						},
					},
					Title: "Question 1",
					Text:  "Block 1",
				},
				{
					State:   2,
					Type:    2,
					Default: 0,
					Options: []dto.Option{},
					Title:   "Question 2",
					Text:    "Block 2",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, value.ErrInvalidNode)
	})

	t.Run("should return an error if bot has invalid selection block", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    2,
					Default: 2, // <----- should be zero
					Options: []dto.Option{
						{
							Text: "To 2",
							Next: 2,
						},
					},
					Title: "Selection 1",
					Text:  "Block 1",
				},
				{
					State:   2,
					Type:    2,
					Default: 0,
					Options: []dto.Option{},
					Title:   "Question 2",
					Text:    "Block 2",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, value.ErrInvalidNode)
	})

	t.Run("should return an error if bot has invalid selection option", func(t *testing.T) {

		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    3,
					Default: 0,
					Options: []dto.Option{
						{
							Text: "To 0",
							Next: 0, // <----- should be non-zero
						},
					},
					Title: "Selection 1",
					Text:  "Block 1",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, value.ErrInvalidOption)
	})

	t.Run("should return an error if bot has invalid block type", func(t *testing.T) {
		ctx := context.Background()

		botRepos := botmemory.NewMemoryBotRepository()
		blcRepos := blockmemory.NewMemoryBlockRepository()

		manager := application.NewBotsManager(
			slogdiscard.NewDiscardLogger(),
			botRepos, blcRepos)

		a := args{
			name:  "example",
			token: "",
			start: 1,
			blocks: []dto.Block{
				{
					State:   1,
					Type:    -1, // <--- should be 1, 2 or 3
					Default: 0,
					Options: []dto.Option{},
					Title:   "Question 1",
					Text:    "Block 1",
				},
			},
		}

		_, err := manager.Create(ctx, a.name, a.token, a.start, a.blocks)
		require.ErrorIs(t, err, value.ErrInvalidNodeType)
	})
}
