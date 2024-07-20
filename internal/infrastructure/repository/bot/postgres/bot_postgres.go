package postgres

import (
	"context"
	"database/sql"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"github.com/pkg/errors"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

type BotPostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresBotRepository(url string) (*BotPostgresRepository, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &BotPostgresRepository{
		db: db,
	}, nil
}

func (r *BotPostgresRepository) Close() error {
	return r.db.Close()
}

func (r *BotPostgresRepository) Save(
	ctx context.Context,
	bot *entity.Bot,
) (value.BotId, error) {
	const op = "PostgresBotsRepository.Save"

	var botId value.BotId
	err := pgutils.Get(
		ctx, r.db, &botId,
		`INSERT INTO bots (Name, Token, Start)
         VALUES ($1, $2, $3)
		 RETURNING Id`,
		bot.Name, bot.Token, bot.Start,
	)
	if err != nil {
		return 0, errors.Wrap(err, op)
	}

	bot.Id = botId

	return botId, nil
}

func (r *BotPostgresRepository) Bot(
	ctx context.Context,
	id value.BotId,
) (*entity.Bot, error) {
	const op = "PostgresBotsRepository.Bot"

	var bot botRow
	err := pgutils.Get(
		ctx, r.db, &bot,
		`SELECT Id, Name, Token, Start
         FROM bots
         WHERE Id = $1`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrBotNotFound
		}
		return nil, errors.Wrap(err, op)
	}

	return botFromRow(bot), nil
}
