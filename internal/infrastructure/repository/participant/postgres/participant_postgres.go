package postgres

import (
	"context"
	"database/sql"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
	"github.com/bmstu-itstech/itsreg-bots/internal/infrastructure/interfaces"
	"github.com/bmstu-itstech/itsreg-bots/pkg/pgerrors"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/zhikh23/pgutils"
)

type ParticipantPostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresParticipantRepository(url string) (*ParticipantPostgresRepository, error) {
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	return &ParticipantPostgresRepository{
		db: db,
	}, nil
}

func (r *ParticipantPostgresRepository) Close() error {
	return r.db.Close()
}

func (r *ParticipantPostgresRepository) Save(
	ctx context.Context,
	prt *entity.Participant,
) error {
	const op = "PostgresParticipantRepository.Save"

	_, err := pgutils.Exec(
		ctx, r.db,
		`INSERT INTO participants (bot_id, user_id, current)
		 VALUES ($1, $2, $3)`,
		prt.Id.BotId, prt.Id.UserId, prt.Current)
	if err != nil {
		if pgerrors.IsUniqueViolationError(err) {
			return interfaces.ErrParticipantAlreadyExists
		}
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *ParticipantPostgresRepository) Participant(
	ctx context.Context,
	id value.ParticipantId,
) (*entity.Participant, error) {
	const op = "PostgresParticipantRepository.Participant"

	var participant participantRow
	err := pgutils.Get(
		ctx, r.db, &participant,
		`SELECT bot_id, user_id, current 
         FROM participants
         WHERE bot_id = $1 AND user_id = $2`,
		id.BotId, id.UserId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrParticipantNotFound
		}
		return nil, errors.Wrap(err, op)
	}

	return participantFromRow(participant), nil
}

func (r *ParticipantPostgresRepository) UpdateState(
	ctx context.Context,
	id value.ParticipantId,
	state value.State,
) error {
	const op = "PostgresParticipantRepository.UpdateState"

	res, err := pgutils.Exec(
		ctx, r.db,
		`UPDATE participants 
         SET current = $1
         WHERE bot_id = $2 AND user_id = $3`,
		state, id.BotId, id.UserId,
	)
	if err != nil {
		return errors.Wrap(err, op)
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, op)
	}

	if aff == 0 {
		return interfaces.ErrParticipantNotFound
	}

	return nil
}
