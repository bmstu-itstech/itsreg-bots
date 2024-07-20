package postgres

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/pkg/pgerrors"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/zhikh23/pgutils"
)

type AnswerPostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresAnswerRepository(url string) (*AnswerPostgresRepository, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &AnswerPostgresRepository{
		db: db,
	}, nil
}

func (r *AnswerPostgresRepository) Close() error {
	return r.db.Close()
}

func (r *AnswerPostgresRepository) Save(
	ctx context.Context,
	ans *entity.Answer,
) error {
	const op = "AnswerPostgresRepository.Save"

	_, err := pgutils.Exec(
		ctx, r.db,
		`INSERT INTO answers (bot_id, user_id, state, value_)
         VALUES ($1, $2, $3, $4)`,
		ans.Id.ParticipantId.BotId, ans.Id.ParticipantId.UserId, ans.Id.State, ans.Value)
	if err != nil {
		if pgerrors.IsUniqueViolationError(err) {
			return interfaces.ErrAnswerExists
		}
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *AnswerPostgresRepository) AnswersFrom(
	ctx context.Context,
	prtId value.ParticipantId,
) ([]*entity.Answer, error) {
	const op = "AnswerPostgresRepository.AnswersFrom"

	var rows []answerRow
	err := pgutils.Select(
		ctx, r.db, &rows,
		`SELECT bot_id, user_id, state, value_
         FROM answers 
         WHERE bot_id = $1 AND user_id = $2`,
		prtId.BotId, prtId.UserId)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return answersFromRows(rows), nil
}
