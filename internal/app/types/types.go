package types

import (
	"time"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type Option struct {
	Text string
	Next int
}

type Block struct {
	Type      string
	State     int
	NextState int
	Options   []Option
	Title     string
	Text      string
}

type EntryPoint struct {
	Key   string
	State int
}

type Mailing struct {
	Name          string
	EntryKey      string
	RequiredState int
}

type Bot struct {
	UUID      string
	Entries   []EntryPoint
	Mailings  []Mailing
	Blocks    []Block
	Name      string
	Token     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AnswersTable struct {
	THead []string
	TBody [][]string
}

func MapOptionFromDomain(option bots.Option) Option {
	return Option{
		Text: option.Text,
		Next: option.Next,
	}
}

func MapOptionToDomain(option Option) (bots.Option, error) {
	return bots.NewOption(option.Text, option.Next)
}

func MapOptionsFromDomain(options []bots.Option) []Option {
	res := make([]Option, len(options))
	for i, option := range options {
		res[i] = MapOptionFromDomain(option)
	}
	return res
}

func MapOptionsToDomain(options []Option) ([]bots.Option, error) {
	res := make([]bots.Option, len(options))
	for i, option := range options {
		o, err := MapOptionToDomain(option)
		if err != nil {
			return nil, err
		}
		res[i] = o
	}
	return res, nil
}

func MapEntryPointFromDomain(entry bots.EntryPoint) EntryPoint {
	return EntryPoint{
		Key:   entry.Key,
		State: entry.State,
	}
}

func MapEntryPointToDomain(entry EntryPoint) (bots.EntryPoint, error) {
	return bots.NewEntryPoint(entry.Key, entry.State)
}

func MapEntriesFromDomain(entries []bots.EntryPoint) []EntryPoint {
	res := make([]EntryPoint, len(entries))
	for i, entry := range entries {
		res[i] = MapEntryPointFromDomain(entry)
	}
	return res
}

func MapEntriesToDomain(entries []EntryPoint) ([]bots.EntryPoint, error) {
	res := make([]bots.EntryPoint, len(entries))
	for i, entry := range entries {
		e, err := MapEntryPointToDomain(entry)
		if err != nil {
			return nil, err
		}
		res[i] = e
	}
	return res, nil
}

func MapMailingFromDomain(mailing bots.Mailing) Mailing {
	return Mailing{
		Name:          mailing.Name,
		EntryKey:      mailing.EntryKey,
		RequiredState: mailing.RequireState,
	}
}

func MapMailingsFromDomain(mailings []bots.Mailing) []Mailing {
	res := make([]Mailing, len(mailings))
	for i, mailing := range mailings {
		res[i] = MapMailingFromDomain(mailing)
	}
	return res
}

func MapMailingToDomain(mailing Mailing) (bots.Mailing, error) {
	return bots.NewMailing(
		mailing.Name,
		mailing.EntryKey,
		mailing.RequiredState,
	)
}

func MapMailingsToDomain(mailings []Mailing) ([]bots.Mailing, error) {
	res := make([]bots.Mailing, len(mailings))
	for i, mailing := range mailings {
		m, err := MapMailingToDomain(mailing)
		if err != nil {
			return nil, err
		}
		res[i] = m
	}
	return res, nil
}

func MapBlockFromDomain(block bots.Block) Block {
	return Block{
		Type:      block.Type.String(),
		State:     block.State,
		NextState: block.NextState,
		Options:   MapOptionsFromDomain(block.Options),
		Title:     block.Title,
		Text:      block.Text,
	}
}

func MapBlockToDomain(block Block) (bots.Block, error) {
	opts, err := MapOptionsToDomain(block.Options)
	if err != nil {
		return bots.Block{}, err
	}
	return bots.NewBlock(
		block.Type,
		block.State,
		block.NextState,
		opts,
		block.Title,
		block.Text,
	)
}

func MapBlocksFromDomain(blocks []bots.Block) []Block {
	res := make([]Block, len(blocks))
	for i, block := range blocks {
		res[i] = MapBlockFromDomain(block)
	}
	return res
}

func MapBlocksToDomain(blocks []Block) ([]bots.Block, error) {
	res := make([]bots.Block, len(blocks))
	for i, block := range blocks {
		b, err := MapBlockToDomain(block)
		if err != nil {
			return nil, err
		}
		res[i] = b
	}
	return res, nil
}

func MapBotFromDomain(bot *bots.Bot) Bot {
	return Bot{
		UUID:      bot.UUID,
		Entries:   MapEntriesFromDomain(bot.Entries()),
		Blocks:    MapBlocksFromDomain(bot.Blocks()),
		Mailings:  MapMailingsFromDomain(bot.Mailings()),
		Name:      bot.Name,
		Token:     bot.Token,
		Status:    bot.Status.String(),
		CreatedAt: bot.CreatedAt,
		UpdatedAt: bot.UpdatedAt,
	}
}

func MapBotsFromDomain(bs []*bots.Bot) []Bot {
	res := make([]Bot, len(bs))
	for i, b := range bs {
		res[i] = MapBotFromDomain(b)
	}
	return res
}

func MapAnswersTableFromDomain(table *bots.AnswersTable) AnswersTable {
	return AnswersTable{
		THead: table.Head,
		TBody: table.Body,
	}
}
