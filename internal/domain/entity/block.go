package entity

import (
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

var (
	ErrInvalidBlockText = errors.New("invalid block text")
)

type Block struct {
	value.Node
	BotId value.BotId
	Title string
	Text  string
}

func NewBlock(node value.Node, botId value.BotId, title string, text string) (*Block, error) {
	if len(text) == 0 {
		return nil, ErrInvalidBlockText
	}

	return &Block{
		Node:  node,
		BotId: botId,
		Title: title,
		Text:  text,
	}, nil
}
