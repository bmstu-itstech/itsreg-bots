package entity

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type Block struct {
	Node  value.Node
	BotId value.BotId
	Title string
	Text  string
}
