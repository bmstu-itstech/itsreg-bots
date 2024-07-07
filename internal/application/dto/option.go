package dto

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type Option struct {
	Text string
	Next uint64
}

func OptionFromDto(dto Option) (value.Option, error) {
	return value.NewOption(dto.Text, value.State(dto.Next))
}

func OptionsFromDto(dto []Option) ([]value.Option, error) {
	opts := make([]value.Option, len(dto))
	for i, dto := range dto {
		opt, err := OptionFromDto(dto)
		if err != nil {
			return nil, err
		}
		opts[i] = opt
	}
	return opts, nil
}
