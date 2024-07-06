package memory

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/interfaces"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	"testing"
)

func TestParticipantMemoryRepository_Save(t *testing.T) {
	t.Run("should save a participant", func(t *testing.T) {
		ctx := context.Background()
		repos := participantMemoryRepository{
			m: map[value.ParticipantId]*entity.Participant{},
		}

		prtId := value.ParticipantId{
			BotId:  42,
			UserId: 123,
		}
		prt, err := entity.NewParticipant(prtId, value.State(3))
		require.NoError(t, err)

		err = repos.Save(ctx, prt)
		require.NoError(t, err)

		got, ok := repos.m[prtId]
		require.True(t, ok)
		require.Equal(t, *prt, *got)
	})

	t.Run("should return error when participant already exists", func(t *testing.T) {
		ctx := context.Background()
		repos := participantMemoryRepository{
			m: map[value.ParticipantId]*entity.Participant{},
		}

		prtId := value.ParticipantId{
			BotId:  42,
			UserId: 123,
		}
		prt, err := entity.NewParticipant(prtId, value.State(3))
		require.NoError(t, err)

		err = repos.Save(ctx, prt)
		require.NoError(t, err)

		prtId = value.ParticipantId{
			BotId:  42,
			UserId: 123,
		}
		prt, err = entity.NewParticipant(prtId, value.State(35))
		require.NoError(t, err)

		err = repos.Save(ctx, prt)
		require.ErrorIs(t, err, interfaces.ErrParticipantAlreadyExists)
	})
}

func TestBlockMemoryRepository_Block(t *testing.T) {
	participants := []*entity.Participant{
		{
			Id: value.ParticipantId{
				BotId:  25,
				UserId: 42,
			},
			Current: value.State(6),
		},
		{
			Id: value.ParticipantId{
				BotId:  25,
				UserId: 37,
			},
			Current: value.State(3),
		},
	}

	repos := participantMemoryRepository{
		m: map[value.ParticipantId]*entity.Participant{},
	}

	for _, prt := range participants {
		repos.m[prt.Id] = prt
	}

	t.Run("should find participant", func(t *testing.T) {
		ctx := context.Background()

		got, err := repos.Participant(ctx, value.ParticipantId{
			BotId:  25,
			UserId: 37,
		})
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, *participants[1], *got)
	})

	t.Run("should return error when participant does not exists", func(t *testing.T) {
		ctx := context.Background()

		_, err := repos.Participant(ctx, value.ParticipantId{
			BotId:  37,
			UserId: 37,
		})
		require.ErrorIs(t, err, interfaces.ErrParticipantNotFound)
	})
}

func TestParticipantMemoryRepository_UpdateState(t *testing.T) {
	participants := []*entity.Participant{
		{
			Id: value.ParticipantId{
				BotId:  25,
				UserId: 42,
			},
			Current: value.State(6),
		},
	}

	repos := participantMemoryRepository{
		m: map[value.ParticipantId]*entity.Participant{},
	}

	for _, prt := range participants {
		repos.m[prt.Id] = prt
	}

	t.Run("should update state", func(t *testing.T) {
		ctx := context.Background()

		prtId := value.ParticipantId{
			BotId:  25,
			UserId: 42,
		}

		newState := value.State(7)
		err := repos.UpdateState(ctx, prtId, newState)
		require.NoError(t, err)

		require.Equal(t, newState, repos.m[prtId].Current)
	})

	t.Run("should return error when participant does not exists", func(t *testing.T) {
		ctx := context.Background()

		prtId := value.ParticipantId{
			BotId:  25,
			UserId: 37,
		}

		newState := value.State(7)
		err := repos.UpdateState(ctx, prtId, newState)
		require.ErrorIs(t, err, interfaces.ErrParticipantNotFound)
	})
}
