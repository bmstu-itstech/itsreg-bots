package dto

import "github.com/zhikh23/itsreg-bots/internal/domain/entity"

type ProcessRequest struct {
	UserId uint64
	BotId  uint64
	Answer string
}

type ProcessResponse struct {
	Messages []Message
}

type Message struct {
	Text    string
	Options []string
}

func MapBlockToMessage(block entity.Block) Message {
	dtoOptions := make([]string, len(block.Node.Options))

	for i, option := range block.Node.Options {
		dtoOptions[i] = option.Text
	}

	return Message{
		Text:    block.Text,
		Options: dtoOptions,
	}
}
