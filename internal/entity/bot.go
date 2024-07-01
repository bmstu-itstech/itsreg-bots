package entity

import "github.com/zhikh23/itsreg-bots/internal/objects"

type Bot struct {
	Id    int64
	Title string
	Token string
	Start objects.NodeId
}
