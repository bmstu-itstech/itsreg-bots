package postgres

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
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
	nonExistentState  = value.State(gofakeit.Uint32())
	nonExistentUserId = value.UserId(gofakeit.Uint32())
)

func setupRepos(t *testing.T) *ParticipantPostgresRepository {
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

	repos, err := NewPostgresParticipantRepository(url)
	require.NoError(t, err)
	require.NotNil(t, repos)

	return repos
}

func TestParticipantPostgresRepository_Save(t *testing.T) {
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

	// Создаем два блока для соблюдения условий внешних ключей.
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO blocks (bot_id, state, type, title, text)
         VALUES ($1, $3, $4, $5, $6), 
                ($2, $3, $4, $5, $6)`,
		botId1, botId2, 1, "message", "test", "test")

	t.Run("should save participant", func(t *testing.T) {
		t.Parallel()

		// Один и тот же пользователь будет участников двух ботов.
		userId := value.UserId(gofakeit.Uint32())

		// Успешно сохраняем пользователя как участника первого бота.
		prt1, err := entity.NewParticipant(
			value.ParticipantId{
				BotId:  botId1,
				UserId: userId,
			}, 1)
		require.NoError(t, err)
		err = repos.Save(ctx, prt1)
		require.NoError(t, err)

		// Успешно сохраняем того же пользователя как участника второго бота.
		prt2, err := entity.NewParticipant(
			value.ParticipantId{
				BotId:  botId2,
				UserId: userId,
			}, 1)
		require.NoError(t, err)
		err = repos.Save(ctx, prt2)
		require.NoError(t, err)

		// Проверяем, что пользователь является участником двух ботов сразу.
		var row participantRow
		err = pgutils.Get(
			ctx, repos.db, &row,
			`SELECT bot_id, user_id, current 
             FROM participants 
             WHERE bot_id = $1 AND user_id = $2`,
			botId1, userId)
		require.NoError(t, err)
		got := participantFromRow(row)
		require.Equal(t, prt1, got)

		err = pgutils.Get(
			ctx, repos.db, &row,
			`SELECT bot_id, user_id, current 
             FROM participants 
             WHERE bot_id = $1 AND user_id = $2`,
			botId2, userId)
		require.NoError(t, err)
		got = participantFromRow(row)
		require.Equal(t, prt2, got)
	})

	t.Run("should return error if participant already exists", func(t *testing.T) {
		t.Parallel()

		// Уникальность участника определяется парой ключей: userId и botId.
		// Проверим тот факт, что невозможно сохранить дважды участника под одной и той
		// же парой ключей.
		userId := value.UserId(gofakeit.Uint32())
		prtId := value.ParticipantId{
			BotId:  botId1,
			UserId: userId,
		}

		// Первый раз участник должен быть успешно сохранен.
		prt1, err := entity.NewParticipant(prtId, 1)
		require.NoError(t, err)
		err = repos.Save(ctx, prt1)
		require.NoError(t, err)

		// Второй раз мы должны получить ошибку ErrParticipantAlreadyExists.
		prt2, err := entity.NewParticipant(prtId, 2)
		require.NoError(t, err)
		err = repos.Save(ctx, prt2)
		require.ErrorIs(t, err, interfaces.ErrParticipantAlreadyExists)
	})

	t.Run("should return error if bot does not exists", func(t *testing.T) {
		t.Parallel()

		// Проверяем наличия условия существования блока.
		userId := value.UserId(gofakeit.Uint32())
		prt, err := entity.NewParticipant(
			value.ParticipantId{
				BotId:  nonExistentBotId,
				UserId: userId,
			}, 1)
		require.NoError(t, err)

		err = repos.Save(ctx, prt)
		require.Error(t, err) // Ошибка не специфицирована.
	})

	t.Run("should return error if state does not exists", func(t *testing.T) {
		t.Parallel()

		// Проверяем наличие условия существования блока,
		// в котором может находится участник.
		userId := value.UserId(gofakeit.Uint32())
		prt, err := entity.NewParticipant(
			value.ParticipantId{
				BotId:  botId1,
				UserId: userId,
			}, nonExistentState)
		require.NoError(t, err)

		err = repos.Save(ctx, prt)
		require.Error(t, err) // Ошибка не специфицирована.
	})
}

func TestParticipantPostgresRepository_Participant(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := setupRepos(t)
	t.Cleanup(func() {
		_ = repos.Close()
	})

	// Создаем двух ботов для соблюдения внешних ключей.
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

	// Добавляем ботам по одному блоку для соблюдения условий внешних ключей.
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO blocks (bot_id, state, type, title, text)
         VALUES ($1, $3, $4, $5, $6), 
                ($2, $3, $4, $5, $6)`,
		botId1, botId2, 1, "message", "test", "test")
	require.NoError(t, err)

	// Добавляем двух поьзователей, причем первый является участником двух ботов сразу.
	userId1 := value.UserId(gofakeit.Uint32())
	userId2 := value.UserId(gofakeit.Uint32())
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO participants (bot_id, user_id, current)
         VALUES ($1, $2, $3),
                ($4, $5, $6),
                ($7, $8, $9)`,
		botId1, userId1, 1,
		botId2, userId1, 1,
		botId2, userId2, 1)
	require.NoError(t, err)

	t.Run("should return a participant", func(t *testing.T) {
		t.Parallel()

		// Должен найти пользователя 1, который является участником бота 2.
		got, err := repos.Participant(ctx, value.ParticipantId{
			BotId:  botId2,
			UserId: userId1,
		})
		require.NoError(t, err)
		require.NotNil(t, got)

		// Должен найти пользователя 2, который является участником бота 2.
		got, err = repos.Participant(ctx, value.ParticipantId{
			BotId:  botId2,
			UserId: userId1,
		})
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("should return error if participant does not exists", func(t *testing.T) {
		t.Parallel()

		// Несуществующий пользователь существующего бота.
		_, err := repos.Participant(ctx, value.ParticipantId{
			BotId:  botId1,
			UserId: nonExistentUserId,
		})
		require.ErrorIs(t, err, interfaces.ErrParticipantNotFound)
	})
}

func TestParticipantPostgresRepository_UpdateState(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := setupRepos(t)
	t.Cleanup(func() {
		_ = repos.Close()
	})

	// Создаем единственного бота для соблюдения условий внешних ключей.
	var botId value.BotId
	err := pgutils.Get(
		ctx, repos.db, &botId,
		`INSERT INTO bots (name, token, start)
         VALUES ('test name', 'test token', 1)
         RETURNING id`,
	)
	require.NoError(t, err)

	// Создаем для бота два блока: один для изначального состояния участников,
	// второй для ожидаемого после обновления.
	// Создание блоков необходимо для соблюдения условий внешних ключей.
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO blocks (bot_id, state, type, title, text)
         VALUES ($1, $2, $4, $5, $6),
                ($1, $3, $4, $5, $6)`,
		botId, 1, 2, "message", "test", "test")
	require.NoError(t, err)

	// Создаем двух участников. Состояние второго мы изменяем,
	// а состояние первого должно остаться константным.
	userId1 := value.UserId(gofakeit.Uint32())
	userId2 := value.UserId(gofakeit.Uint32())
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO participants (bot_id, user_id, current)
         VALUES ($1, $2, $4),
                ($1, $3, $4)`,
		botId, userId1, userId2, 1)
	require.NoError(t, err)

	t.Run("should update participant state", func(t *testing.T) {
		t.Parallel()

		expected := value.State(2) // Ожидаем сменить состояние на 2.
		err := repos.UpdateState(ctx, value.ParticipantId{
			BotId:  botId,
			UserId: userId2,
		}, expected)
		require.NoError(t, err)

		// Проверяем, что состояние требуемого участника изменилось на ожидаемое.
		var got value.State
		err = pgutils.Get(
			ctx, repos.db, &got,
			`SELECT current FROM participants WHERE bot_id = $1 AND user_id = $2`,
			botId, userId2)
		require.NoError(t, err)
		require.Equal(t, expected, got)

		// Проверяем, что не изменилось состояние другого участника.
		err = pgutils.Get(
			ctx, repos.db, &got,
			`SELECT current FROM participants WHERE bot_id = $1 AND user_id = $2`,
			botId, userId1)
		require.NoError(t, err)
		require.NotEqual(t, expected, got)
	})

	t.Run("should return error if participant does not exists", func(t *testing.T) {
		t.Parallel()

		expected := value.State(2) // Ожидаем сменить состояние на 2.
		err := repos.UpdateState(ctx, value.ParticipantId{
			BotId:  botId,
			UserId: nonExistentUserId,
		}, expected)
		require.ErrorIs(t, err, interfaces.ErrParticipantNotFound)
	})
}
