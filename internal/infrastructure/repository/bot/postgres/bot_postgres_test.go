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
	nonExistentBotId = value.BotId(gofakeit.Uint32())
)

func setupRepos(t *testing.T) *BotPostgresRepository {
	cfg := config.MustLoadPath(configPath)
	url := endpoint.BuildPostgresConnectionString(
		endpoint.WithPostgresUser(cfg.Postgres.User),
		endpoint.WithPostgresPassword(cfg.Postgres.Pass),
		endpoint.WithPostgresHost(cfg.Postgres.Host),
		endpoint.WithPostgresPort(cfg.Postgres.Port),
		endpoint.WithPostgresDb(cfg.Postgres.DbName),
	)

	repos, err := NewPostgresBotRepository(url)
	require.NoError(t, err)
	require.NotNil(t, repos)

	return repos
}

func TestBotPostgresRepository_Save(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := setupRepos(t)
	t.Cleanup(func() {
		_ = repos.Close()
	})

	t.Run("should save a new bot", func(t *testing.T) {
		t.Parallel()

		// Сохраняем бота.
		bot, err := entity.NewBot(value.UnknownBotId, "test", "test", 1)
		require.NoError(t, err)
		id, err := repos.Save(ctx, bot)
		require.NoError(t, err)
		require.NotZero(t, id)
		require.Equal(t, bot.Id, id)

		// Проверяем, что бот был сохранен в БД.
		var row botRow
		err = pgutils.Get(
			ctx, repos.db, &row,
			`SELECT Id, Name, Token, Start
             FROM bots WHERE Id = $1`,
			id,
		)
		require.NoError(t, err)
		got := botFromRow(row)
		require.Equal(t, bot, got)
	})
}

func TestBotPostgresRepository_Bot(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := setupRepos(t)
	t.Cleanup(func() {
		_ = repos.Close()
	})

	// Сохраняем тестового бота.
	bot, err := entity.NewBot(value.UnknownBotId, "test", "test", 1)
	require.NoError(t, err)
	var botId value.BotId
	err = pgutils.Get(
		ctx, repos.db, &botId,
		`INSERT INTO bots (name, token, start)
	     VALUES ($1, $2, $3)
         RETURNING id`,
		bot.Name, bot.Token, bot.Start,
	)
	require.NoError(t, err)
	bot.Id = botId // Чтобы далее можно было провести полное сравнение по полям,
	// обновляем Id бота с тем, что сейчас в БД.

	t.Run("should get a bot", func(t *testing.T) {
		t.Parallel()
		got, err := repos.Bot(ctx, botId)
		require.NoError(t, err)
		require.Equal(t, got, bot)
	})

	t.Run("should return error if bot does not exist", func(t *testing.T) {
		t.Parallel()
		_, err := repos.Bot(ctx, nonExistentBotId)
		require.ErrorIs(t, err, interfaces.ErrBotNotFound)
	})
}
