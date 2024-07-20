package postgres

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/config"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/pkg/endpoint"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
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
	nonExistentState = value.State(gofakeit.Uint32())
)

func setupRepos(t *testing.T) *BlockPostgresRepository {
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

	repos, err := NewPostgresBlockRepository(url)
	require.NoError(t, err)
	require.NotNil(t, repos)

	return repos
}

func TestBlockPostgresRepository_Save(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := setupRepos(t)
	t.Cleanup(func() {
		_ = repos.Close()
	})

	// Создаем ботов для соблюдения условия внешнего ключа.
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

	t.Run("should save block without options", func(t *testing.T) {
		t.Parallel()

		state := value.State(gofakeit.Uint32())
		node, err := value.NewQuestionNode(state, value.StateNone)
		require.NoError(t, err)
		block, err := entity.NewBlock(node, botId1, "test", "test")
		require.NoError(t, err)

		err = repos.Save(ctx, block)
		require.NoError(t, err)

		// Проверем, что блок реально был сохранен в БД.
		var row blockRow
		err = pgutils.Get(
			ctx, repos.db, &row,
			`SELECT bot_id, state, type, default_, Title, Text
			 FROM blocks
			 WHERE bot_id = $1 AND state = $2`,
			block.BotId, node.State,
		)
		require.NoError(t, err)
		got := blockFromRow(row)
		require.Equal(t, block, got)
	})

	t.Run("should save block with options", func(t *testing.T) {
		t.Parallel()

		state := value.State(gofakeit.Uint32())
		node, err := value.NewSelectionNode(state, []value.Option{
			{Text: "option 1", Next: state},
			{Text: "option 2", Next: 0},
		})
		require.NoError(t, err)
		block, err := entity.NewBlock(node, botId1, "test", "test")
		require.NoError(t, err)

		err = repos.Save(ctx, block)
		require.NoError(t, err)

		// Проверяем, что блок реально был сохранен в БД.
		var row blockRow
		err = pgutils.Get(
			ctx, repos.db, &row,
			`SELECT bot_id, state, type, default_, Title, Text
			 FROM blocks
			 WHERE bot_id = $1 AND state = $2`,
			block.BotId, node.State,
		)
		require.NoError(t, err)
		got := blockFromRow(row)

		// Полное сравнение по полям как
		// > require.Equal(t, block, got)
		// невозможно, так запрос возвращает блок без опций.
		require.Equal(t, block.BotId, got.BotId)
		require.Equal(t, block.State, got.State)
		require.Equal(t, block.Default, got.Default)
		require.Equal(t, block.Title, got.Title)
		require.Equal(t, block.Text, got.Text)

		// Проверяем, что опции тоже были сохранены.
		var rows []optionRow
		err = pgutils.Select(
			ctx, repos.db, &rows,
			`SELECT bot_id, state, Next, value_
             FROM options
             WHERE bot_id = $1 AND State = $2`,
			botId1, node.State,
		)
		require.NoError(t, err)
		require.Len(t, rows, len(node.Options))
		options := optionsFromRows(rows)
		require.ElementsMatch(t, options, node.Options) // Порядок не важен.
	})

	t.Run("should save blocks with same states in different bots", func(t *testing.T) {
		t.Parallel()

		// Одно состояние для двух блоков в разных ботах.
		state := value.State(gofakeit.Uint32())

		// Блок в первом боте.
		node, err := value.NewQuestionNode(state, value.StateNone)
		require.NoError(t, err)
		block1, err := entity.NewBlock(node, botId1, "test", "test")
		require.NoError(t, err)
		err = repos.Save(ctx, block1)
		require.NoError(t, err)

		// Блок во втором боте.
		node, err = value.NewQuestionNode(state, value.StateNone)
		require.NoError(t, err)
		block2, err := entity.NewBlock(node, botId2, "test", "test")
		require.NoError(t, err)
		err = repos.Save(ctx, block2)
		require.NoError(t, err)

		// Проверяем, что блоки были сохранены в БД.
		var row blockRow
		err = pgutils.Get(
			ctx, repos.db, &row,
			`SELECT bot_id, state, type, default_, Title, Text
			 FROM blocks
			 WHERE bot_id = $1 AND state = $2`,
			block1.BotId, node.State,
		)
		require.NoError(t, err)
		got := blockFromRow(row)
		require.Equal(t, block1, got)

		err = pgutils.Get(
			ctx, repos.db, &row,
			`SELECT bot_id, state, type, default_, Title, Text
			 FROM blocks
			 WHERE bot_id = $1 AND state = $2`,
			block2.BotId, node.State,
		)
		require.NoError(t, err)
		got = blockFromRow(row)
		require.Equal(t, block2, got)
	})

	t.Run("should return error if state already exists", func(t *testing.T) {
		t.Parallel()

		// Одно состояние для двух блоков в одном боте.
		state := value.State(gofakeit.Uint32())

		// Создаем два блока в одном боте с одним и тем же состоянием.
		// Репозиторий должен вернуть ошибку.

		// Состояние хранится в Узле, поэтому его можно создать единожды.
		node, err := value.NewQuestionNode(state, value.StateNone)
		require.NoError(t, err)

		// Первый блок должен быть сохранен без ошибки.
		block1, err := entity.NewBlock(node, botId1, "test", "test")
		require.NoError(t, err)
		err = repos.Save(ctx, block1)
		require.NoError(t, err)

		// Для второго блока для того же бота с тем же состоянием
		// репозиторий должен вернуть ошибку BlockAlreadyExists.
		block2, err := entity.NewBlock(node, botId1, "test", "test")
		require.NoError(t, err)
		err = repos.Save(ctx, block2)
		require.ErrorIs(t, err, interfaces.ErrBlockAlreadyExists)
	})

	t.Run("should return error if bot id not exists", func(t *testing.T) {
		t.Parallel()

		// Проверяем случай, когда блок имеет ссылку на несуществующего бота.
		state := value.State(gofakeit.Uint32())
		node, err := value.NewMessageNode(state, value.StateNone)
		require.NoError(t, err)
		block, err := entity.NewBlock(node, nonExistentBotId, "test", "test")
		require.NoError(t, err)

		err = repos.Save(ctx, block)
		require.Error(t, err) // Ошибка не специфицирована.
	})
}

func TestBlockPostgresRepository_Block(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	repos := setupRepos(t)
	t.Cleanup(func() {
		_ = repos.Close()
	})

	// Создаем бота для соблюдения условий внешних ключей.
	var botId value.BotId
	err := pgutils.Get(
		ctx, repos.db, &botId,
		`INSERT INTO bots (name, token, start)
         VALUES ('test name', 'test token', 1)
         RETURNING id`,
	)
	require.NoError(t, err)

	// Блок без опций.
	state1 := value.State(gofakeit.Uint32())
	node1, err := value.NewQuestionNode(state1, value.StateNone)
	require.NoError(t, err)
	block1, err := entity.NewBlock(node1, botId, "test", "test")
	require.NoError(t, err)
	_, err = pgutils.Exec(
		ctx, repos.db,
		`INSERT INTO blocks (bot_id, state, type, default_, title, text)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		block1.BotId, block1.State, nodeTypeToString(block1.Type),
		block1.Default, block1.Title, block1.Text,
	)
	require.NoError(t, err)

	// Блок с опциями.
	state2 := value.State(gofakeit.Uint32())
	options := []value.Option{
		{Text: "option 1", Next: state1},
		{Text: "option 2", Next: state2},
	}
	node2, err := value.NewSelectionNode(state2, options)
	require.NoError(t, err)
	block2, err := entity.NewBlock(node2, botId, "test", "test")
	require.NoError(t, err)

	err = pgutils.RunTx(ctx, repos.db, func(tx *sqlx.Tx) error {
		_, err = pgutils.Exec(
			ctx, tx,
			`INSERT INTO blocks (bot_id, state, type, default_, title, text)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			block2.BotId, block2.State, nodeTypeToString(block2.Type),
			block2.Default, block2.Title, block2.Text,
		)
		if err != nil {
			return err
		}

		for _, option := range options {
			_, err = pgutils.Exec(
				ctx, tx,
				`INSERT INTO options (bot_id, state, Next, value_)
				 VALUES ($1, $2, $3, $4)`,
				block2.BotId, block2.State, option.Next, option.Text,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
	require.NoError(t, err)

	t.Run("should return block without options", func(t *testing.T) {
		got, err := repos.Block(ctx, botId, state1)
		require.NoError(t, err)
		require.Equal(t, block1, got)
	})

	t.Run("should return block with options", func(t *testing.T) {
		got, err := repos.Block(ctx, botId, state2)
		require.NoError(t, err)

		// Сравнение через
		// > require.Equal(t, block2, got)
		// не подойдет так как порядок опций не определен.
		require.Equal(t, block2.BotId, got.BotId)
		require.Equal(t, block2.State, got.State)
		require.Equal(t, block2.Default, got.Default)
		require.Equal(t, block2.Title, got.Title)
		require.Equal(t, block2.Text, got.Text)
		require.ElementsMatch(t, options, got.Options)
	})

	t.Run("should return error if block does not exists", func(t *testing.T) {
		_, err := repos.Block(ctx, botId, nonExistentState)
		require.ErrorIs(t, err, interfaces.ErrBlockNotFound)
	})
}
