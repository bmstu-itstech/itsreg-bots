package postgres

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/pkg/endpoint"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/zhikh23/pgutils"
	"testing"
)

const (
	// Для запуска этих тестов требуется запущенное локальное окружение (docker-compose.local).
	configPath = "../../../../../config/local.yaml"
)

var (
	nonExistentBotId  = value.BotId(gofakeit.Uint32())
	nonExistentUserId = value.UserId(gofakeit.Uint32())
)

func setupRepos(t *testing.T) *AnswerPostgresRepository {
	// Чтобы гарантировать успешно подключение к тестовой БД,
	// используем те же настройки, что использует приложение
	// в локальной конфигурации.
	cfg := config.MustLoadPath(configPath)
	url := endpoint.BuildPostgresConnectionString(
		endpoint.WithPostgresUser(cfg.Postgres.User),
		endpoint.WithPostgresPassword(cfg.Postgres.Pass),
		endpoint.WithPostgresHost(cfg.Postgres.Host),
		endpoint.WithPostgresPort(cfg.Postgres.Port),
		endpoint.WithPostgresDb(cfg.Postgres.DbName),
	)

	repos, err := NewPostgresAnswerRepository(url)
	require.NoError(t, err)
	require.NotNil(t, repos)

	return repos
}

func TestAnswerPostgresRepository_Save(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := setupRepos(t)
	t.Cleanup(func() {
		_ = repos.Close()
	})

	// Создаем двух ботов, в которых один и тот же пользователь будет
	// участником одновременно.
	var botId1 value.BotId
	err := pgutils.Get(
		ctx, repos.db, &botId1,
		`INSERT INTO bots (name, token, start)
         VALUES ('test name', 'test token', 1)
         RETURNING id`,
	)
	require.NoError(t, err)

	var botId2 value.BotId
	err = pgutils.Get(
		ctx, repos.db, &botId2,
		`INSERT INTO bots (name, token, start)
         VALUES ('test name', 'test token', 1)
         RETURNING id`,
	)
	require.NoError(t, err)

	// Три блока будут иметь два уникальных состояния.
	// Так два блока в разных ботах будут иметь одинаковые состояния.
	state1 := value.State(gofakeit.Uint32())
	state2 := value.State(gofakeit.Uint32())

	// Создаем три блока для соблюдения условий внешних ключей:
	// При этом выполняются условия:
	// - существуют два блока в одном боте, но с разными состояниями;
	// - существуют два блока с одинаковыми состояниями, но в разных ботах;
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO blocks (bot_id, state, type, title, text)
         VALUES ($1, $2, $3, $4, $5),
                ($6, $7, $8, $9, $10),
                ($11, $12, $13, $14, $15)`,
		botId1, state1, "question", "test", "test",
		botId1, state2, "question", "test", "test",
		botId2, state2, "question", "test", "test",
	)

	// Создаем участников для соблюдения условий внешних ключей.
	// При этом выполняются условия:
	// - существуют пользователь, который является участником двух ботов;
	// - существует бот, который имеет несколько участников.
	userId1 := value.UserId(gofakeit.Uint32())
	userId2 := value.UserId(gofakeit.Uint32())
	prtId1 := value.ParticipantId{
		BotId:  botId1,
		UserId: userId1,
	}
	prtId2 := value.ParticipantId{
		BotId:  botId1,
		UserId: userId2,
	}
	prtId3 := value.ParticipantId{
		BotId:  botId2,
		UserId: userId2,
	}
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO participants (bot_id, user_id, current)
         VALUES ($1, $2, $3),
                ($4, $5, $6),
                ($7, $8, $9)`,
		prtId1.BotId, prtId1.UserId, state1,
		prtId2.BotId, prtId2.UserId, state1,
		prtId3.BotId, prtId3.UserId, state2,
	)
	require.NoError(t, err)

	t.Run("should save answers", func(t *testing.T) {
		t.Parallel()

		// Сохраняем несколько ответов.
		// При этом выполняются условия:
		// - каждый участник имеет хотя бы один ответ;
		// - существует участник с более чем одним ответом.
		ans1, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: prtId1,
				State:         state1,
			}, "answer 1")
		require.NoError(t, err)
		err = repos.Save(ctx, ans1)
		require.NoError(t, err)

		ans2, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: prtId1,
				State:         state2,
			}, "answer 2")
		require.NoError(t, err)
		err = repos.Save(ctx, ans2)
		require.NoError(t, err)

		ans3, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: prtId2,
				State:         state2,
			}, "answer 3")
		require.NoError(t, err)
		err = repos.Save(ctx, ans3)
		require.NoError(t, err)

		ans4, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: prtId3,
				State:         state2,
			}, "answer 4")
		require.NoError(t, err)
		err = repos.Save(ctx, ans4)
		require.NoError(t, err)

		// Проверяем, что все ответы сохранились.
		var got string

		err = pgutils.Get(
			ctx, repos.db, &got,
			`SELECT value_
             FROM answers
             WHERE bot_id = $1 AND user_id = $2 AND state = $3`,
			ans1.Id.ParticipantId.BotId, ans1.Id.ParticipantId.UserId, ans1.Id.State,
		)
		require.NoError(t, err)
		require.Equal(t, ans1.Value, got)

		err = pgutils.Get(
			ctx, repos.db, &got,
			`SELECT value_
             FROM answers
             WHERE bot_id = $1 AND user_id = $2 AND state = $3`,
			ans2.Id.ParticipantId.BotId, ans2.Id.ParticipantId.UserId, ans2.Id.State,
		)
		require.NoError(t, err)
		require.Equal(t, ans2.Value, got)

		err = pgutils.Get(
			ctx, repos.db, &got,
			`SELECT value_
             FROM answers
             WHERE bot_id = $1 AND user_id = $2 AND state = $3`,
			ans3.Id.ParticipantId.BotId, ans3.Id.ParticipantId.UserId, ans3.Id.State,
		)
		require.NoError(t, err)
		require.Equal(t, ans3.Value, got)

		err = pgutils.Get(
			ctx, repos.db, &got,
			`SELECT value_
             FROM answers
             WHERE bot_id = $1 AND user_id = $2 AND state = $3`,
			ans4.Id.ParticipantId.BotId, ans4.Id.ParticipantId.UserId, ans4.Id.State,
		)
		require.NoError(t, err)
		require.Equal(t, ans4.Value, got)
	})

	t.Run("should return error if bot does not exists", func(t *testing.T) {
		// Попытка сохранить ответ для несуществующего бота.
		ans, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: value.ParticipantId{
					BotId:  nonExistentBotId,
					UserId: userId1,
				},
				State: state1,
			}, "answer")
		require.NoError(t, err)
		err = repos.Save(ctx, ans)
		require.Error(t, err) // Ошибка не специфицирована.
	})

	t.Run("should return error if state does not exists", func(t *testing.T) {
		// Попытка сохранить ответ для несуществующего состояния бота.
		ans, err := entity.NewAnswer(
			value.AnswerId{
				ParticipantId: value.ParticipantId{
					BotId:  botId2,
					UserId: userId2,
				},
				State: state1, // Бот №2 не имеет состояния state1.
			}, "answer")
		require.NoError(t, err)
		err = repos.Save(ctx, ans)
		require.Error(t, err) // Ошибка не специфицирована.
	})
}

func TestAnswerPostgresRepository_AnswersFrom(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := setupRepos(t)
	t.Cleanup(func() {
		_ = repos.Close()
	})

	// Создаем двух ботов, в которых один и тот же пользователь будет
	// участником одновременно.
	var botId1 value.BotId
	err := pgutils.Get(
		ctx, repos.db, &botId1,
		`INSERT INTO bots (name, token, start)
         VALUES ('test name', 'test token', 1)
         RETURNING id`,
	)
	require.NoError(t, err)

	var botId2 value.BotId
	err = pgutils.Get(
		ctx, repos.db, &botId2,
		`INSERT INTO bots (name, token, start)
         VALUES ('test name', 'test token', 1)
         RETURNING id`,
	)
	require.NoError(t, err)

	// Три блока будут иметь два уникальных состояния.
	// Так два блока в разных ботах будут иметь одинаковые состояния.
	state1 := value.State(gofakeit.Uint32())
	state2 := value.State(gofakeit.Uint32())

	// Создаем три блока для соблюдения условий внешних ключей:
	// При этом выполняются условия:
	// - существуют два блока в одном боте, но с разными состояниями;
	// - существуют два блока с одинаковыми состояниями, но в разных ботах;
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO blocks (bot_id, state, type, title, text)
         VALUES ($1, $2, $3, $4, $5),
                ($6, $7, $8, $9, $10),
                ($11, $12, $13, $14, $15)`,
		botId1, state1, "question", "test", "test",
		botId1, state2, "question", "test", "test",
		botId2, state2, "question", "test", "test",
	)

	// Создаем участников для соблюдения условий внешних ключей.
	// При этом выполняются условия:
	// - существуют пользователь, который является участником двух ботов;
	// - существует бот, который имеет несколько участников.
	userId1 := value.UserId(gofakeit.Uint32())
	userId2 := value.UserId(gofakeit.Uint32())
	prtId1 := value.ParticipantId{
		BotId:  botId1,
		UserId: userId1,
	}
	prtId2 := value.ParticipantId{
		BotId:  botId1,
		UserId: userId2,
	}
	prtId3 := value.ParticipantId{
		BotId:  botId2,
		UserId: userId2,
	}
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO participants (bot_id, user_id, current)
         VALUES ($1, $2, $3),
                ($4, $5, $6),
                ($7, $8, $9)`,
		prtId1.BotId, prtId1.UserId, state1,
		prtId2.BotId, prtId2.UserId, state1,
		prtId3.BotId, prtId3.UserId, state2,
	)
	require.NoError(t, err)

	// Сохраняем несколько ответов.
	// При этом выполняются условия:
	// - каждый участник имеет хотя бы один ответ;
	// - существует участник с более чем одним ответом.
	ans1, err := entity.NewAnswer(
		value.AnswerId{
			ParticipantId: prtId1,
			State:         state1,
		}, "answer 1")
	require.NoError(t, err)
	err = repos.Save(ctx, ans1)
	require.NoError(t, err)

	ans2, err := entity.NewAnswer(
		value.AnswerId{
			ParticipantId: prtId1,
			State:         state2,
		}, "answer 2")
	require.NoError(t, err)
	err = repos.Save(ctx, ans2)
	require.NoError(t, err)

	ans3, err := entity.NewAnswer(
		value.AnswerId{
			ParticipantId: prtId2,
			State:         state2,
		}, "answer 3")
	require.NoError(t, err)
	err = repos.Save(ctx, ans3)
	require.NoError(t, err)

	ans4, err := entity.NewAnswer(
		value.AnswerId{
			ParticipantId: prtId3,
			State:         state2,
		}, "answer 4")
	require.NoError(t, err)
	err = repos.Save(ctx, ans4)
	require.NoError(t, err)

	t.Run("should return one answer", func(t *testing.T) {
		res, err := repos.AnswersFrom(ctx, prtId2)
		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, ans3, res[0])
	})

	t.Run("should return multiple answers", func(t *testing.T) {
		res, err := repos.AnswersFrom(ctx, prtId1)
		require.NoError(t, err)
		require.Len(t, res, 2)
		require.ElementsMatch(t, []*entity.Answer{ans1, ans2}, res)
	})

	t.Run("should return empty answers", func(t *testing.T) {
		res, err := repos.AnswersFrom(ctx, value.ParticipantId{
			BotId:  botId1,
			UserId: nonExistentUserId,
		})
		require.NoError(t, err)
		require.Empty(t, res)
	})
}
