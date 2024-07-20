package postgres

import (
	"context"
	"database/sql"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/pkg/pgerrors"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/zhikh23/pgutils"
)

type BlockPostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresBlockRepository(url string) (*BlockPostgresRepository, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &BlockPostgresRepository{
		db: db,
	}, nil
}

func (r *BlockPostgresRepository) Close() error {
	return r.db.Close()
}

func (r *BlockPostgresRepository) Save(
	ctx context.Context,
	block *entity.Block,
) error {
	if len(block.Options) > 0 {
		return r.saveBlockWithOptions(ctx, block)
	}

	return r.saveBlockWithoutOptions(ctx, block)
}

func (r *BlockPostgresRepository) saveBlockWithoutOptions(
	ctx context.Context,
	block *entity.Block,
) error {
	const op = "BlockPostgresRepository.Save"

	_, err := pgutils.Exec(
		ctx, r.db,
		`INSERT INTO blocks (bot_id, state, type, default_, title, text) 
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		block.BotId, block.State, nodeTypeToString(block.Type), block.Default, block.Title, block.Text,
	)
	if err != nil {
		if pgerrors.IsUniqueViolationError(err) {
			return interfaces.ErrBlockAlreadyExists
		}
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *BlockPostgresRepository) saveBlockWithOptions(
	ctx context.Context,
	block *entity.Block,
) error {
	const op = "BlockPostgresRepository.Save"

	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		_, err := pgutils.Exec(
			ctx, tx,
			`INSERT INTO blocks (bot_id, state, type, default_, title, text) 
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			block.BotId, block.State, nodeTypeToString(block.Type), block.Default, block.Title, block.Text,
		)
		if err != nil {
			return err
		}

		for _, opt := range block.Options {
			_, err = pgutils.Exec(
				ctx, tx,
				`INSERT INTO options (bot_id, state, next, value_)
                 VALUES ($1, $2, $3, $4)`,
				block.BotId, block.State, opt.Next, opt.Text,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *BlockPostgresRepository) Block(
	ctx context.Context,
	botId value.BotId,
	state value.State,
) (*entity.Block, error) {
	const op = "BlockPostgresRepository.Block"

	var row blockRow
	err := pgutils.Get(
		ctx, r.db, &row,
		`SELECT bot_id, state, type, default_, title, text
	     FROM blocks
		 WHERE bot_id = $1 AND state = $2
		 LIMIT 1`,
		botId, state,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrBlockNotFound
		}
		return nil, errors.Wrap(err, op)
	}
	block := blockFromRow(row)

	// Блок может иметь опции.
	var rows []optionRow
	err = pgutils.Select(
		ctx, r.db, &rows,
		`SELECT next, value_
         FROM options
         WHERE bot_id = $1 AND state = $2`,
		botId, state,
	)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}
	block.Options = optionsFromRows(rows)

	return block, nil
}
