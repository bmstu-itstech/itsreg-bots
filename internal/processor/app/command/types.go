package command

import "github.com/bmstu-itstech/itsreg-bots/internal/processor/domain/bots"

type Option struct {
	Text string
	Next int
}

func mapOptionsToDomain(options []Option) ([]bots.Option, error) {
	res := make([]bots.Option, len(options))
	for i, option := range options {
		opt, err := bots.NewOption(option.Text, option.Next)
		if err != nil {
			return nil, err
		}
		res[i] = opt
	}
	return res, nil
}

type Block struct {
	Type      string
	State     int
	NextState int
	Options   []Option
	Title     string
	Text      string
}

func mapBlocksToDomain(blocks []Block) ([]bots.Block, error) {
	res := make([]bots.Block, len(blocks))
	for i, block := range blocks {
		opts, err := mapOptionsToDomain(block.Options)
		if err != nil {
			return nil, err
		}
		b, err := bots.NewBlock(
			block.Type,
			block.State,
			block.NextState,
			opts,
			block.Title,
			block.Text,
		)
		if err != nil {
			return nil, err
		}
		res[i] = b
	}
	return res, nil
}
