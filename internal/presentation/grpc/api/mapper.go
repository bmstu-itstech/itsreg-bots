package api

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/application/dto"
	botsv1 "github.com/zhikh23/itsreg-proto/gen/go/bots"
)

func messageFromDto(msg dto.Message) *botsv1.Message {
	options := make([]*botsv1.MessageOption, len(msg.Options))

	for i, opt := range msg.Options {
		options[i] = &botsv1.MessageOption{
			Text: opt,
		}
	}

	return &botsv1.Message{
		Text:    msg.Text,
		Options: options,
	}
}

func messagesFromDtos(msgs []dto.Message) []*botsv1.Message {
	res := make([]*botsv1.Message, len(msgs))

	for i, m := range msgs {
		res[i] = messageFromDto(m)
	}

	return res
}

func blockToDto(block *botsv1.Block) dto.Block {
	return dto.Block{
		State:   block.State,
		Type:    int(block.Type),
		Default: block.Default,
		Options: optionsToDtos(block.Options),
		Title:   block.Title,
		Text:    block.Text,
	}
}

func blocksToDtos(blocks []*botsv1.Block) []dto.Block {
	dtos := make([]dto.Block, len(blocks))
	for i, block := range blocks {
		dtos[i] = blockToDto(block)
	}
	return dtos
}

func optionToDto(option *botsv1.BlockOption) dto.Option {
	return dto.Option{
		Text: option.Text,
		Next: option.Next,
	}
}

func optionsToDtos(options []*botsv1.BlockOption) []dto.Option {
	dtos := make([]dto.Option, len(options))
	for i, option := range options {
		dtos[i] = optionToDto(option)
	}
	return dtos
}
