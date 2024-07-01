package entity

import "github.com/zhikh23/itsreg-bots/internal/objects"

type Module struct {
	BotId int64
	Title string
	Text  string
	objects.Node
}
