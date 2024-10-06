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

func (r *pgBotsRepository) UpdateOrCreate(ctx context.Context, bot *bots.Bot) error {
	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `DELETE FROM bots WHERE uuid = $1`, bot.UUID)
		if err != nil {
			return err
		}

		if err = r.checkExecRes(tx.NamedExecContext(ctx,
			`INSERT INTO bots 
				(uuid, name, token, status, created_at, updated_at, owner_uuid)
             VALUES (:uuid, :name, :token, :status, :created_at, :updated_at, :owner_uuid)`,
			convertBotToDB(bot),
		)); err != nil {
			return err
		}

		if err = r.checkExecRes(tx.NamedExecContext(ctx,
			`INSERT INTO blocks
				(bot_uuid, state, type, next_state, title, text) 
			 VALUES (:bot_uuid, :state, :type, :next_state, :title, :text)`,
			convertBlocksToDB(bot.UUID, bot.Blocks()),
		)); err != nil {
			return err
		}

		for _, block := range bot.Blocks() {
			if len(block.Options) == 0 {
				continue
			}
			if err = r.checkExecRes(tx.NamedExecContext(ctx,
				`INSERT INTO options
					(bot_uuid, state, next, text) 
				 VALUES (:bot_uuid, :state, :next, :text)`,
				convertOptionsToDB(bot.UUID, block.State, block.Options),
			)); err != nil {
				return err
			}
		}

		if err = r.checkExecRes(tx.NamedExecContext(ctx,
			`INSERT INTO entry_points 
				(bot_uuid, key, state) 
			 VALUES (:bot_uuid, :key, :state)`,
			convertEntryPointsToDB(bot.UUID, bot.Entries()),
		)); err != nil {
			return err
		}

		if len(bot.Mailings()) > 0 {
			if err = r.checkExecRes(tx.NamedExecContext(ctx,
				`INSERT INTO mailings 
				(bot_uuid, name, entry_key, required_state) 
			 VALUES (:bot_uuid, :name, :entry_key, :required_state)`,
				convertMailingsToDB(bot.UUID, bot.Mailings()),
			)); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *pgBotsRepository) UpdateStatus(ctx context.Context, botUUID string, status bots.Status) error {
	err := r.checkExecRes(r.db.ExecContext(ctx,
		`UPDATE bots SET status = $1 WHERE uuid = $2`, status.String(), botUUID,
	))
	if errors.Is(err, ErrNoAffectedRows) {
		return bots.BotNotFoundError{UUID: botUUID}
	}
	return err
}

func (r *pgBotsRepository) Delete(ctx context.Context, uuid string) error {
	err := r.checkExecRes(r.db.ExecContext(ctx, `DELETE FROM bots WHERE uuid = $1`, uuid))
	if errors.Is(err, ErrNoAffectedRows) {
		return bots.BotNotFoundError{UUID: uuid}
	}
	return err
}

func (r *pgBotsRepository) Bot(ctx context.Context, uuid string) (*bots.Bot, error) {
	var bRow botRow
	if err := pgutils.Get(ctx, r.db, &bRow,
		`SELECT uuid, name, token, status, created_at, updated_at, owner_uuid
         FROM   bots 
		 WHERE  uuid = $1`, uuid,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, bots.BotNotFoundError{UUID: uuid}
	} else if err != nil {
		return nil, err
	}

	entryPoints, err := r.selectEntryPoints(ctx, bRow.UUID)
	if err != nil {
		return nil, err
	}

	mailings, err := r.selectMailings(ctx, bRow.UUID)
	if err != nil {
		return nil, err
	}

	blocks, err := r.selectBlocks(ctx, bRow.UUID)
	if err != nil {
		return nil, err
	}

	return bots.UnmarshallBotFromDB(
		bRow.UUID, bRow.OwnerUUID, entryPoints, mailings, blocks,
		bRow.Name, bRow.Token, bRow.Status,
		bRow.CreatedAt.Local(), bRow.UpdatedAt.Local(),
	)
}

func (r *pgBotsRepository) UserBots(ctx context.Context, userUUID string) ([]*bots.Bot, error) {
	var bRows []botRow
	if err := pgutils.Select(ctx, r.db, &bRows,
		`SELECT uuid, name, token, status, created_at, updated_at, owner_uuid
         FROM   bots 
		 WHERE  owner_uuid = $1`, userUUID,
	); err != nil {
		return nil, err
	}

	res := make([]*bots.Bot, 0, len(bRows))
	for _, bRow := range bRows {
		entryPoints, err := r.selectEntryPoints(ctx, bRow.UUID)
		if err != nil {
			return nil, err
		}

		mailings, err := r.selectMailings(ctx, bRow.UUID)
		if err != nil {
			return nil, err
		}

		blocks, err := r.selectBlocks(ctx, bRow.UUID)
		if err != nil {
			return nil, err
		}

		bot, err := bots.UnmarshallBotFromDB(
			bRow.UUID, bRow.OwnerUUID, entryPoints, mailings, blocks,
			bRow.Name, bRow.Token, bRow.Status,
			bRow.CreatedAt.Local(), bRow.UpdatedAt.Local(),
		)
		if err != nil {
			return nil, err
		}

		res = append(res, bot)
	}

	return res, nil
}

func (r *pgBotsRepository) selectEntryPoints(ctx context.Context, uuid string) ([]bots.EntryPoint, error) {
	var eRows []entryPointRow
	if err := pgutils.Select(ctx, r.db, &eRows,
		`SELECT bot_uuid, key, state 
		 FROM   entry_points 
         WHERE  bot_uuid = $1`, uuid,
	); err != nil {
		return nil, err
	}
	return convertEntryPointsToDomain(eRows)
}

func (r *pgBotsRepository) selectMailings(ctx context.Context, uuid string) ([]bots.Mailing, error) {
	var mRows []mailingRow
	if err := pgutils.Select(ctx, r.db, &mRows,
		`SELECT bot_uuid, name, entry_key, required_state 
		 FROM   mailings 
         WHERE  bot_uuid = $1`, uuid,
	); err != nil {
		return nil, err
	}
	return convertMailingsToDomain(mRows)
}

func (r *pgBotsRepository) selectOptions(ctx context.Context, uuid string, state int) ([]bots.Option, error) {
	var oRows []optionRow
	if err := pgutils.Select(ctx, r.db, &oRows,
		`SELECT bot_uuid, state, text, next 
		 FROM   options 
		 WHERE  bot_uuid = $1 AND state = $2`, uuid, state,
	); err != nil {
		return nil, err
	}
	return convertOptionsToDomain(oRows)
}

func (r *pgBotsRepository) selectBlocks(ctx context.Context, uuid string) ([]bots.Block, error) {
	var bRows []blockRow
	if err := pgutils.Select(ctx, r.db, &bRows,
		`SELECT bot_uuid, state, type, next_state, title, text
		 FROM   blocks
         WHERE  bot_uuid = $1`, uuid,
	); err != nil {
		return nil, err
	}

	blocks := make([]bots.Block, 0, len(bRows))
	for _, row := range bRows {
		options, err := r.selectOptions(ctx, uuid, row.State)
		if err != nil {
			return nil, err
		}

		block, err := convertBlockToDomain(row, options)
		if err != nil {
			return nil, err
		}

		blocks = append(blocks, block)
	}

	return blocks, nil
}

var ErrNoAffectedRows = errors.New("no affected rows")

func (r *pgBotsRepository) checkExecRes(res sql.Result, err error) error {
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return ErrNoAffectedRows
	}

	return nil
}

type optionRow struct {
	BotUUID string `db:"bot_uuid"`
	State   int    `db:"state"`
	Text    string `db:"text"`
	Next    *int   `db:"next"`
}

func convertOptionToDB(botUUID string, state int, o bots.Option) optionRow {
	return optionRow{
		BotUUID: botUUID,
		State:   state,
		Text:    o.Text,
		Next:    nilOnZero(o.Next),
	}
}

func convertOptionsToDB(botUUID string, state int, os []bots.Option) []optionRow {
	res := make([]optionRow, len(os))
	for i, o := range os {
		res[i] = convertOptionToDB(botUUID, state, o)
	}
	return res
}

func convertOptionsToDomain(os []optionRow) ([]bots.Option, error) {
	res := make([]bots.Option, len(os))
	for i, o := range os {
		option, err := bots.NewOption(o.Text, zeroOnNil(o.Next))
		if err != nil {
			return nil, err
		}
		res[i] = option
	}
	return res, nil
}

type entryPointRow struct {
	BotUUID string `db:"bot_uuid"`
	Key     string `db:"key"`
	State   int    `db:"state"`
}

func convertEntryPointToDB(botUUID string, e bots.EntryPoint) entryPointRow {
	return entryPointRow{
		BotUUID: botUUID,
		Key:     e.Key,
		State:   e.State,
	}
}

func convertEntryPointsToDB(botUUID string, es []bots.EntryPoint) []entryPointRow {
	res := make([]entryPointRow, len(es))
	for i, e := range es {
		res[i] = convertEntryPointToDB(botUUID, e)
	}
	return res
}

func convertEntryPointsToDomain(es []entryPointRow) ([]bots.EntryPoint, error) {
	res := make([]bots.EntryPoint, len(es))
	for i, e := range es {
		entryPoint, err := bots.NewEntryPoint(e.Key, e.State)
		if err != nil {
			return nil, err
		}
		res[i] = entryPoint
	}
	return res, nil
}

type mailingRow struct {
	BotUUID       string `db:"bot_uuid"`
	Name          string `db:"name"`
	EntryKey      string `db:"entry_key"`
	RequiredState *int   `db:"required_state"`
}

func convertMailingToDB(botUUID string, m bots.Mailing) mailingRow {
	return mailingRow{
		BotUUID:       botUUID,
		Name:          m.Name,
		EntryKey:      m.EntryKey,
		RequiredState: nilOnZero(m.RequireState),
	}
}

func convertMailingsToDB(botUUID string, ms []bots.Mailing) []mailingRow {
	res := make([]mailingRow, len(ms))
	for i, m := range ms {
		res[i] = convertMailingToDB(botUUID, m)
	}
	return res
}

func convertMailingsToDomain(ms []mailingRow) ([]bots.Mailing, error) {
	res := make([]bots.Mailing, len(ms))
	for i, m := range ms {
		mailing, err := bots.NewMailing(m.Name, m.EntryKey, zeroOnNil(m.RequiredState))
		if err != nil {
			return nil, err
		}
		res[i] = mailing
	}
	return res, nil
}

type blockRow struct {
	BotUUID   string `db:"bot_uuid"`
	Type      string `db:"type"`
	State     int    `db:"state"`
	NextState *int   `db:"next_state"`
	Title     string `db:"title"`
	Text      string `db:"text"`
}

func nilOnZero(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func zeroOnNil(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func convertBlockToDB(botUUID string, b bots.Block) blockRow {
	return blockRow{
		BotUUID:   botUUID,
		Type:      b.Type.String(),
		State:     b.State,
		NextState: nilOnZero(b.NextState),
		Title:     b.Title,
		Text:      b.Text,
	}
}

func convertBlocksToDB(botUUID string, bs []bots.Block) []blockRow {
	res := make([]blockRow, len(bs))
	for i, b := range bs {
		res[i] = convertBlockToDB(botUUID, b)
	}
	return res
}

func convertBlockToDomain(b blockRow, options []bots.Option) (bots.Block, error) {
	return bots.UnmarshallBlockFromDB(b.Type, b.State, zeroOnNil(b.NextState), options, b.Title, b.Text)
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

func convertBotToDB(b *bots.Bot) botRow {
	return botRow{
		UUID:      b.UUID,
		OwnerUUID: b.OwnerUUID,
		Name:      b.Name,
		Token:     b.Token,
		Status:    b.Status.String(),
		CreatedAt: b.CreatedAt.UTC(),
		UpdatedAt: b.UpdatedAt.UTC(),
	}
}
