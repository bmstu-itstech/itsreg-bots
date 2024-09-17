package infra

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type pgBotsRepository struct {
	db *sqlx.DB
}

func NewPgBotsRepository(db *sqlx.DB) bots.Repository {
	return &pgBotsRepository{
		db: db,
	}
}

func (r *pgBotsRepository) Save(ctx context.Context, bot *bots.Bot) error {
	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		if err := insertBot(ctx, tx, bot); err != nil {
			return err
		}

		for _, block := range bot.Blocks() {
			if err := insertBlock(ctx, tx, bot.UUID, block); err != nil {
				return err
			}
			for _, option := range block.Options {
				if err := insertOption(ctx, tx, bot.UUID, block.State, option); err != nil {
					return err
				}
			}
		}

		for _, entry := range bot.Entries() {
			if err := insertEntryPoint(ctx, tx, bot.UUID, entry); err != nil {
				return err
			}
		}

		for _, mailing := range bot.Mailings() {
			if err := insertMailing(ctx, tx, bot.UUID, mailing); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *pgBotsRepository) Bot(ctx context.Context, uuid string) (*bots.Bot, error) {
	return selectBot(ctx, r.db, uuid)
}

func (r *pgBotsRepository) Bots(ctx context.Context, userUUID string) ([]*bots.Bot, error) {
	return selectUserBots(ctx, r.db, userUUID)
}

func (r *pgBotsRepository) Update(
	ctx context.Context,
	uuid string,
	updateFn func(context.Context, *bots.Bot) error,
) error {
	bot, err := selectBot(ctx, r.db, uuid)
	if err != nil {
		return err
	}

	err = updateFn(ctx, bot)
	if err != nil {
		return err
	}

	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		if err = updateBot(ctx, tx, bot); err != nil {
			return err
		}

		if err = updateBlocks(ctx, r.db, uuid, bot.Blocks()); err != nil {
			return err
		}

		for _, block := range bot.Blocks() {
			if err = updateOptions(ctx, r.db, uuid, block.State, block.Options); err != nil {
				return err
			}
		}

		if err = updateEntryPoints(ctx, r.db, uuid, bot.Entries()); err != nil {
			return err
		}

		if err = updateMailings(ctx, r.db, uuid, bot.Mailings()); err != nil {
			return err
		}

		return nil
	})
}

func (r *pgBotsRepository) Delete(ctx context.Context, uuid string) error {
	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		return deleteBot(ctx, tx, uuid)
	})
}

func checkInsertResult(res sql.Result) error {
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return errors.New("0 rows affected")
	}

	return nil
}

func insertBot(ctx context.Context, ex sqlx.ExecerContext, bot *bots.Bot) error {
	res, err := pgutils.Exec(ctx, ex,
		`INSERT INTO
        	bots (uuid, owner_uuid, name, token, status, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		bot.UUID, bot.OwnerUUID, bot.Name, bot.Token, bot.Status.String(),
		bot.CreatedAt.UTC(), bot.UpdatedAt.UTC(),
	)
	if pgutils.IsUniqueViolationError(err) {
		return bots.BotAlreadyExistsError{UUID: bot.UUID}
	} else if err != nil {
		return err
	}

	return checkInsertResult(res)
}

func selectBot(ctx context.Context, q sqlx.QueryerContext, uuid string) (*bots.Bot, error) {
	var row botRow
	err := pgutils.Get(ctx, q, &row,
		`SELECT uuid, owner_uuid, name, token, status, created_at, updated_at
		 FROM   bots
		 WHERE  uuid = $1`, uuid,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, bots.BotNotFoundError{UUID: uuid}
	} else if err != nil {
		return nil, err
	}

	blocks, err := selectBlocks(ctx, q, row.UUID)
	if err != nil {
		return nil, err
	}

	mailings, err := selectMailings(ctx, q, row.UUID)
	if err != nil {
		return nil, err
	}

	entries, err := selectEntryPoints(ctx, q, row.UUID)
	if err != nil {
		return nil, err
	}

	return bots.UnmarshallBotFromDB(
		row.UUID, row.OwnerUUID, entries, mailings, blocks, row.Name, row.Token,
		row.Status, row.CreatedAt.Local(), row.UpdatedAt.Local(),
	)
}

func selectUserBots(ctx context.Context, q sqlx.QueryerContext, userUUID string) ([]*bots.Bot, error) {
	var rows []botRow
	err := pgutils.Select(ctx, q, &rows,
		`SELECT uuid, owner_uuid, name, token, status, created_at, updated_at
		 FROM   bots
		 WHERE  owner_uuid = $1`, userUUID,
	)
	if err != nil {
		return nil, err
	}

	res := make([]*bots.Bot, 0, len(rows))
	for _, row := range rows {
		blocks, err := selectBlocks(ctx, q, row.UUID)
		if err != nil {
			return nil, err
		}

		entries, err := selectEntryPoints(ctx, q, row.UUID)
		if err != nil {
			return nil, err
		}

		mailings, err := selectMailings(ctx, q, row.UUID)
		if err != nil {
			return nil, err
		}

		bot, err := bots.UnmarshallBotFromDB(
			row.UUID, row.OwnerUUID, entries, mailings, blocks, row.Name, row.Token,
			row.Status, row.CreatedAt.Local(), row.UpdatedAt.Local(),
		)
		if err != nil {
			return nil, err
		}

		res = append(res, bot)
	}

	return res, nil
}

func updateBot(ctx context.Context, ex sqlx.ExecerContext, bot *bots.Bot) error {
	res, err := pgutils.Exec(ctx, ex,
		`UPDATE bots 
		 SET    name = $2,
                token = $3,
                status = $4,
                updated_at = $5
		 WHERE  uuid = $1`,
		bot.UUID, bot.Name, bot.Token, bot.Status.String(), bot.UpdatedAt.UTC(),
	)
	if err != nil {
		return err
	}

	return checkInsertResult(res)
}

func deleteBot(ctx context.Context, ex sqlx.ExecerContext, uuid string) error {
	res, err := pgutils.Exec(ctx, ex, `DELETE FROM bots WHERE uuid = $1`, uuid)
	if err != nil {
		return err
	}

	return checkInsertResult(res)
}

func insertEntryPoint(ctx context.Context, ex sqlx.ExecerContext, botUUID string, entry bots.EntryPoint) error {
	res, err := pgutils.Exec(ctx, ex,
		`INSERT INTO
			entry_points (bot_uuid, key, state) 
        VALUES ($1, $2, $3)`,
		botUUID, entry.Key, entry.State,
	)
	if err != nil {
		return err
	}
	return checkInsertResult(res)
}

func updateEntryPoints(
	ctx context.Context, ex sqlx.ExecerContext, botUUID string, entries []bots.EntryPoint,
) error {
	if err := deleteEntryPoints(ctx, ex, botUUID); err != nil {
		return err
	}

	for _, entry := range entries {
		if err := insertEntryPoint(ctx, ex, botUUID, entry); err != nil {
			return err
		}
	}

	return nil
}

func deleteEntryPoints(ctx context.Context, ex sqlx.ExecerContext, botUUID string) error {
	_, err := pgutils.Exec(ctx, ex, `DELETE FROM entry_points WHERE bot_uuid = $1`, botUUID)
	if err != nil {
		return err
	}

	return nil
}

func insertMailing(ctx context.Context, ex sqlx.ExecerContext, botUUID string, mailing bots.Mailing) error {
	res, err := pgutils.Exec(ctx, ex,
		`INSERT INTO mailings
			(bot_uuid, name, entry_key, required_state) 
		 VALUES ($1, $2, $3, $4)`,
		botUUID, mailing.Name, mailing.EntryKey, mailing.RequireState,
	)
	if err != nil {
		return err
	}

	return checkInsertResult(res)
}

func deleteMailings(ctx context.Context, ex sqlx.ExecerContext, botUUID string) error {
	_, err := pgutils.Exec(ctx, ex, `DELETE FROM mailings WHERE bot_uuid = $1`, botUUID)
	if err != nil {
		return err
	}

	return nil
}

func updateMailings(ctx context.Context, ex sqlx.ExecerContext, botUUID string, mailings []bots.Mailing) error {
	if err := deleteMailings(ctx, ex, botUUID); err != nil {
		return err
	}

	for _, mailing := range mailings {
		if err := insertMailing(ctx, ex, botUUID, mailing); err != nil {
			return err
		}
	}

	return nil
}

func insertOption(
	ctx context.Context, ex sqlx.ExecerContext, botUUID string, blockState int, option bots.Option,
) error {
	res, err := pgutils.Exec(ctx, ex,
		`INSERT INTO 
			options (bot_uuid, state, next, text)
		 VALUES ($1, $2, $3, $4)`,
		botUUID, blockState, option.Next, option.Text,
	)
	if err != nil {
		return err
	}
	return checkInsertResult(res)
}

func updateOptions(
	ctx context.Context, ex sqlx.ExecerContext, botUUID string, blockState int, options []bots.Option,
) error {
	if err := deleteOptions(ctx, ex, botUUID, blockState); err != nil {
		return err
	}

	for _, option := range options {
		if err := insertOption(ctx, ex, botUUID, blockState, option); err != nil {
			return err
		}
	}

	return nil
}

func deleteOptions(
	ctx context.Context, ex sqlx.ExecerContext, botUUID string, blockState int,
) error {
	_, err := pgutils.Exec(ctx, ex, `DELETE FROM options WHERE bot_uuid = $1 AND state = $2`, botUUID, blockState)
	if err != nil {
		return err
	}

	return nil
}

func insertBlock(ctx context.Context, ex sqlx.ExecerContext, botUUID string, block bots.Block) error {
	res, err := pgutils.Exec(ctx, ex,
		`INSERT INTO 
			blocks (bot_uuid, state, type, next_state, title, text)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		botUUID, block.State, block.Type.String(), block.NextState, block.Title, block.Text,
	)
	if err != nil {
		return err
	}
	return checkInsertResult(res)
}

func updateBlocks(ctx context.Context, ex sqlx.ExecerContext, botUUID string, blocks []bots.Block) error {
	if err := deleteBlocks(ctx, ex, botUUID); err != nil {
		return err
	}

	for _, block := range blocks {
		if err := insertBlock(ctx, ex, botUUID, block); err != nil {
			return err
		}
	}

	return nil
}

func deleteBlocks(ctx context.Context, ex sqlx.ExecerContext, botUUID string) error {
	res, err := pgutils.Exec(ctx, ex, `DELETE FROM blocks WHERE bot_uuid = $1`, botUUID)
	if err != nil {
		return err
	}

	return checkInsertResult(res)
}

func selectOptions(
	ctx context.Context, q sqlx.QueryerContext, botUUID string, blockState int,
) ([]bots.Option, error) {
	var rows []optionRow
	err := pgutils.Select(ctx, q, &rows,
		`SELECT next, text FROM options WHERE bot_uuid = $1 AND state = $2`,
		botUUID, blockState,
	)
	if err != nil {
		return nil, err
	}

	return convertOptionsFromDB(rows)
}

func convertOptionsFromDB(rows []optionRow) ([]bots.Option, error) {
	res := make([]bots.Option, len(rows))
	for i, row := range rows {
		o, err := bots.NewOption(row.Text, row.Next)
		if err != nil {
			return nil, err
		}
		res[i] = o
	}
	return res, nil
}

func selectEntryPoints(
	ctx context.Context, q sqlx.QueryerContext, botUUID string,
) ([]bots.EntryPoint, error) {
	var rows []entryPointRow
	err := pgutils.Select(ctx, q, &rows,
		`SELECT key, state
		 FROM   entry_points
		 WHERE  bot_uuid = $1`, botUUID,
	)
	if err != nil {
		return nil, err
	}

	return convertEntryPointsFromDB(rows)
}

func convertEntryPointsFromDB(rows []entryPointRow) ([]bots.EntryPoint, error) {
	res := make([]bots.EntryPoint, len(rows))
	for i, row := range rows {
		e, err := bots.NewEntryPoint(row.Key, row.State)
		if err != nil {
			return nil, err
		}
		res[i] = e
	}
	return res, nil
}

func selectMailings(
	ctx context.Context, q sqlx.QueryerContext, botUUID string,
) ([]bots.Mailing, error) {
	var rows []mailingRow
	err := pgutils.Select(ctx, q, &rows,
		`SELECT name, entry_key, required_state 
         FROM   mailings 
         WHERE  bot_uuid = $1`, botUUID,
	)
	if err != nil {
		return nil, err
	}

	return convertMailingsFromDB(rows)
}

func convertMailingsFromDB(rows []mailingRow) ([]bots.Mailing, error) {
	res := make([]bots.Mailing, len(rows))
	for i, row := range rows {
		m, err := bots.NewMailing(row.Name, row.EntryKey, row.RequiredState)
		if err != nil {
			return nil, err
		}
		res[i] = m
	}
	return res, nil
}

func selectBlocks(
	ctx context.Context, q sqlx.QueryerContext, botUUID string,
) ([]bots.Block, error) {
	var rows []blockRow
	err := pgutils.Select(ctx, q, &rows,
		`SELECT type, state, next_state, title, text
		 FROM   blocks
		 WHERE  bot_uuid = $1`, botUUID,
	)
	if err != nil {
		return nil, err
	}

	res := make([]bots.Block, len(rows))
	for i, row := range rows {
		options, err := selectOptions(ctx, q, botUUID, row.State)
		if err != nil {
			return nil, err
		}
		block, err := bots.UnmarshallBlockFromDB(row.Type, row.State, row.NextState, options, row.Title, row.Text)
		if err != nil {
			return nil, err
		}
		res[i] = block
	}

	return res, nil
}

type optionRow struct {
	Text string `db:"text"`
	Next int    `db:"next"`
}

type entryPointRow struct {
	Key   string `db:"key"`
	State int    `db:"state"`
}

type mailingRow struct {
	Name          string `db:"name"`
	EntryKey      string `db:"entry_key"`
	RequiredState int    `db:"required_state"`
}

type blockRow struct {
	Type      string `db:"type"`
	State     int    `db:"state"`
	NextState int    `db:"next_state"`
	Title     string `db:"title"`
	Text      string `db:"text"`
}

type botRow struct {
	UUID      string    `db:"uuid"`
	OwnerUUID string    `db:"owner_uuid"`
	Name      string    `db:"name"`
	Token     string    `db:"token"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
