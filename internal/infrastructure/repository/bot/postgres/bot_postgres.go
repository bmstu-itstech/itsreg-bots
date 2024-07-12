package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

type botPostgresRepository struct {
	db *sqlx.DB
}

type RepositoryOption func(context.Context, *botPostgresRepository) error

func NewPostgresBotRepository(url string) (interfaces.BotRepository, error) {
	fmt.Println(url)
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &botPostgresRepository{
		db: db,
	}, nil
}

func (r *botPostgresRepository) Close() error {
	return r.db.Close()
}

func (r *botPostgresRepository) Save(
	ctx context.Context,
	bot *entity.Bot,
) (value.BotId, error) {
	const op = "PostgresBotsRepository.Save"

	stmt, err := r.db.Prepare(
		"INSERT INTO bots (name, token, start) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}

	var id int64
	err = stmt.QueryRowContext(ctx, bot.Name, bot.Token, bot.Start).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return value.BotId(id), nil
}

func (r *botPostgresRepository) Bot(
	ctx context.Context,
	id value.BotId,
) (*entity.Bot, error) {
	const op = "PostgresBotsRepository.Bot"

	stmt, err := r.db.Prepare("SELECT id, name, token, start FROM bots WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var bot entity.Bot
	err = row.Scan(&bot.Id, &bot.Name, &bot.Token, &bot.Start)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrBotNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &bot, nil
}
