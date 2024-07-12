package memory

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAnswerMemoryRepository_Save(t *testing.T) {
	t.Run("should save an answer", func(t *testing.T) {
		ctx := context.Background()
		repos := answerMemoryRepository{
			m: map[value.AnswerId]*entity.Answer{},
		}

		ans, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: value.ParticipantId{
					BotId:  value.BotId(2),
					UserId: value.UserId(61),
				},
				State: value.State(24),
			}, "some user's answer")
		require.NoError(t, err)

		err = repos.Save(ctx, ans)
		require.NoError(t, err)

		got, ok := repos.m[ans.Id]
		require.True(t, ok)
		require.Equal(t, *ans, *got)
	})

	t.Run("should return error when answer already exists", func(t *testing.T) {
		ctx := context.Background()
		repos := answerMemoryRepository{
			m: map[value.AnswerId]*entity.Answer{},
		}

		ans, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: value.ParticipantId{
					BotId:  value.BotId(2),
					UserId: value.UserId(61),
				},
				State: value.State(24),
			}, "some user's answer")
		require.NoError(t, err)

		err = repos.Save(ctx, ans)
		require.NoError(t, err)

		another, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: value.ParticipantId{
					BotId:  value.BotId(2),
					UserId: value.UserId(61),
				},
				State: value.State(24),
			}, "some another user's answer with same id")
		require.NoError(t, err)

		err = repos.Save(ctx, another)
		require.ErrorIs(t, err, interfaces.ErrAnswerAlreadyExists)
	})
}

func TestAnswerMemoryRepository_AnswersFrom(t *testing.T) {
	answers := []*entity.Answer{
		{
			Id: value.AnswerId{
				ParticipantId: value.ParticipantId{
					BotId:  value.BotId(2),
					UserId: value.UserId(61),
				},
				State: value.State(24),
			},
			Value: "some user's answer",
		},
		{
			Id: value.AnswerId{
				ParticipantId: value.ParticipantId{
					BotId:  value.BotId(2),
					UserId: value.UserId(61),
				},
				State: value.State(25),
			},
			Value: "some another user's answer",
		},
		{
			Id: value.AnswerId{
				ParticipantId: value.ParticipantId{
					BotId:  value.BotId(2),
					UserId: value.UserId(37),
				},
				State: value.State(24),
			},
			Value: "some another user",
		},
	}

	repos := answerMemoryRepository{m: map[value.AnswerId]*entity.Answer{}}
	for _, ans := range answers {
		repos.m[ans.Id] = ans
	}

	t.Run("should find answers", func(t *testing.T) {
		ctx := context.Background()

		prtId := value.ParticipantId{
			BotId:  value.BotId(2),
			UserId: value.UserId(61),
		}

		got, err := repos.AnswersFrom(ctx, prtId)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Len(t, got, 2)

		for _, got := range got {
			require.Equal(t, got.Id.ParticipantId, prtId)
		}
	})

	t.Run("should return empty answers", func(t *testing.T) {
		ctx := context.Background()

		prtId := value.ParticipantId{
			BotId:  value.BotId(2),
			UserId: value.UserId(7),
		}

		got, err := repos.AnswersFrom(ctx, prtId)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Empty(t, got)
	})
}
