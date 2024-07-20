package postgres

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

type botRow struct {
	Id    int64
	Name  string
	Token string
	Start int64
}

func botFromRow(row botRow) *entity.Bot {
	return &entity.Bot{
		Id:    value.BotId(row.Id),
		Name:  row.Name,
		Token: row.Token,
		Start: value.State(row.Start),
	}
}
