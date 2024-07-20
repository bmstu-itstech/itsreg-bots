package postgres

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

func nodeTypeToString(typ value.NodeType) string {
	switch typ {
	case value.Message:
		return "message"
	case value.Question:
		return "question"
	case value.Selection:
		return "selection"
	}
	return "unknown"
}

func nodeTypeFromString(s string) value.NodeType {
	switch s {
	case "message":
		return value.Message
	case "question":
		return value.Question
	case "selection":
		return value.Selection
	}
	return 0
}

type blockRow struct {
	BotId   int64  `db:"bot_id"`
	State   int64  `db:"state"`
	Type    string `db:"type"`
	Default int64  `db:"default_"`
	Title   string `db:"title"`
	Text    string `db:"text"`
}

func blockFromRow(row blockRow) *entity.Block {
	return &entity.Block{
		Node: value.Node{
			Type:    nodeTypeFromString(row.Type),
			State:   value.State(row.State),
			Default: value.State(row.Default),
			Options: make([]value.Option, 0),
		},
		BotId: value.BotId(row.BotId),
		Title: row.Title,
		Text:  row.Text,
	}
}

type optionRow struct {
	BotId int64  `db:"bot_id"`
	State int64  `db:"state"`
	Value string `db:"value_"`
	Next  int64  `db:"next"`
}

func optionFromRow(row optionRow) value.Option {
	return value.Option{
		Text: row.Value,
		Next: value.State(row.Next),
	}
}

func optionsFromRows(rows []optionRow) []value.Option {
	res := make([]value.Option, len(rows))
	for i, row := range rows {
		res[i] = optionFromRow(row)
	}
	return res
}
