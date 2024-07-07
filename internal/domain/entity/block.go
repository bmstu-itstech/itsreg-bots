package entity

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

var (
	ErrInvalidBlock = errors.New("invalid block")
)

type Block struct {
	Node  value.Node
	BotId value.BotId
	Title string
	Text  string
}

func NewBlock(node value.Node, botId value.BotId, title string, text string) (*Block, error) {
	if len(text) == 0 {
		return nil, ErrInvalidBlock
	}

	return &Block{
		Node:  node,
		BotId: botId,
		Title: title,
		Text:  text,
	}, nil
}
