package infra

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type pgParticipantsRepository struct {
	db *sqlx.DB
}

func NewPgParticipantsRepository(db *sqlx.DB) bots.ParticipantRepository {
	return &pgParticipantsRepository{
		db: db,
	}
}

func (r *pgParticipantsRepository) ParticipantsOfBot(ctx context.Context, botUUID string) ([]*bots.Participant, error) {
	return selectParticipants(ctx, r.db, botUUID)
}

func (r *pgParticipantsRepository) UpdateOrCreate(
	ctx context.Context,
	botUUID string,
	userID int64,
	updateFn func(context.Context, *bots.Participant) error,
) error {
	prt, err := selectParticipant(ctx, r.db, botUUID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		prt, err = bots.NewParticipant(botUUID, userID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = updateFn(ctx, prt)
	if err != nil {
		return err
	}

	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err = upsertParticipant(ctx, tx, prt)
		if err != nil {
			return err
		}

		for _, ans := range prt.Answers() {
			err = upsertAnswer(ctx, tx, botUUID, userID, ans)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func selectParticipants(
	ctx context.Context, q sqlx.QueryerContext, botUUID string,
) ([]*bots.Participant, error) {
	var rows []participantRow
	err := pgutils.Select(ctx, q, &rows,
		`SELECT
			bot_uuid, user_id, state
		FROM
			participants
		WHERE 
			bot_uuid = $1`, botUUID,
	)
	if err != nil {
		return nil, err
	}

	res := make([]*bots.Participant, len(rows))
	for i, row := range rows {
		answers, err := selectAnswers(ctx, q, row.BotUUID, row.UserID)
		if err != nil {
			return nil, err
		}
		prt, err := bots.UnmarshallParticipantFromDB(row.BotUUID, row.UserID, row.State, answers)
		if err != nil {
			return nil, err
		}
		res[i] = prt
	}

	return res, nil
}

func selectParticipant(
	ctx context.Context, q sqlx.QueryerContext, botUUID string, userID int64,
) (*bots.Participant, error) {
	var row participantRow
	err := pgutils.Get(ctx, q, &row,
		`SELECT 
			bot_uuid, user_id, state
		FROM 
			participants 
		WHERE 
			bot_uuid = $1 AND user_id = $2`,
		botUUID, userID,
	)
	if err != nil {
		return nil, err
	}

	answers, err := selectAnswers(ctx, q, botUUID, userID)
	if err != nil {
		return nil, err
	}

	return bots.UnmarshallParticipantFromDB(row.BotUUID, row.UserID, row.State, answers)
}

func upsertParticipant(ctx context.Context, ex sqlx.ExecerContext, prt *bots.Participant) error {
	res, err := pgutils.Exec(ctx, ex,
		`INSERT INTO
			participants (bot_uuid, user_id, state)
		VALUES 
			($1, $2, $3)
		ON CONFLICT ( bot_uuid, user_id )
			DO UPDATE SET state = EXCLUDED.state`,
		prt.BotUUID, prt.UserID, prt.State,
	)
	if err != nil {
		return err
	}

	return checkInsertResult(res)
}

func selectAnswers(
	ctx context.Context, q sqlx.QueryerContext, botUUID string, userID int64,
) ([]bots.Answer, error) {
	var rows []answerRow
	err := pgutils.Select(ctx, q, &rows,
		`SELECT
			state, text
		FROM 
			answers
		WHERE bot_uuid = $1 AND user_id = $2`,
		botUUID, userID,
	)
	if err != nil {
		return nil, err
	}

	return mapAnswersFromDB(rows)
}

func upsertAnswer(ctx context.Context, ex sqlx.ExecerContext, botUUID string, userID int64, ans bots.Answer) error {
	res, err := pgutils.Exec(ctx, ex,
		`INSERT INTO
			answers (bot_uuid, user_id, state, text)
		VALUES
			($1, $2, $3, $4)
		ON CONFLICT ( bot_uuid, user_id, state )
			DO UPDATE SET text = EXCLUDED.text`,
		botUUID, userID, ans.State, ans.Text,
	)
	if err != nil {
		return err
	}

	return checkInsertResult(res)
}

func mapAnswersFromDB(rows []answerRow) ([]bots.Answer, error) {
	res := make([]bots.Answer, len(rows))
	for i, row := range rows {
		a, err := bots.NewAnswer(row.State, row.Text)
		if err != nil {
			return nil, err
		}
		res[i] = a
	}
	return res, nil
}

func checkInsertResult(res sql.Result) error {
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return errors.New("no affected rows")
	}

	return nil
}

type participantRow struct {
	BotUUID string `db:"bot_uuid"`
	UserID  int64  `db:"user_id"`
	State   int    `db:"state"`
}

type answerRow struct {
	State int    `db:"state"`
	Text  string `db:"text"`
}
