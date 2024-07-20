package postgres

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

type participantRow struct {
	BotId   int64 `db:"bot_id"`
	UserId  int64 `db:"user_id"`
	Current int64 `db:"current"`
}

func participantFromRow(row participantRow) *entity.Participant {
	return &entity.Participant{
		Id: value.ParticipantId{
			BotId:  value.BotId(row.BotId),
			UserId: value.UserId(row.UserId),
		},
		Current: value.State(row.Current),
	}
}
