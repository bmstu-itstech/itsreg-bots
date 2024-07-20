package postgres

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

type answerRow struct {
	BotId  int64  `db:"bot_id"`
	UserId int64  `db:"user_id"`
	State  int64  `db:"state"`
	Value  string `db:"value_"`
}

func answerFromRow(row answerRow) *entity.Answer {
	return &entity.Answer{
		Id: value.AnswerId{
			ParticipantId: value.ParticipantId{
				BotId:  value.BotId(row.BotId),
				UserId: value.UserId(row.UserId),
			},
			State: value.State(row.State),
		},
		Value: row.Value,
	}
}

func answersFromRows(rows []answerRow) []*entity.Answer {
	res := make([]*entity.Answer, len(rows))
	for i, row := range rows {
		res[i] = answerFromRow(row)
	}
	return res
}
