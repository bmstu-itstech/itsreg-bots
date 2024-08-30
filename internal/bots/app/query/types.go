package query

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
	"time"
)

type Table struct {
	Head []string
	Body [][]string
}

func mapTableFromDomain(table *bots.Table) Table {
	return Table{
		Head: table.Head,
		Body: table.Body,
	}
}

type Option struct {
	Text string
	Next int
}

func mapOptionsFromDomain(options []bots.Option) []Option {
	res := make([]Option, len(options))
	for i, option := range options {
		opt := Option{Text: option.Text, Next: option.Next}
		res[i] = opt
	}
	return res
}

type Block struct {
	Type      string
	State     int
	NextState int
	Options   []Option
	Title     string
	Text      string
}

func mapBlocksFromDomain(blocks map[int]bots.Block) []Block {
	res := make([]Block, 0, len(blocks))
	for _, block := range blocks {
		opts := mapOptionsFromDomain(block.Options)
		b := Block{
			Type:      block.Type.String(),
			State:     block.State,
			NextState: block.NextState,
			Options:   opts,
			Title:     block.Title,
			Text:      block.Text,
		}
		res = append(res, b)
	}
	return res
}

type Bot struct {
	UUID       string
	Status     string
	Blocks     []Block
	StartState int
	Name       string
	Token      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func mapBotFromDomain(bot *bots.Bot) Bot {
	return Bot{
		UUID:       bot.UUID,
		Status:     bot.Status.String(),
		Blocks:     mapBlocksFromDomain(bot.Blocks),
		StartState: bot.StartState,
		Name:       bot.Name,
		Token:      bot.Token,
		CreatedAt:  bot.CreatedAt,
		UpdatedAt:  bot.UpdatedAt,
	}
}
