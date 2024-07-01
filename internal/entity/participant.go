package entity

import "github.com/zhikh23/itsreg-bots/internal/objects"

type ParticipantId struct {
	BotId  int64
	UserId int64
}

type Participant struct {
	CurrentId objects.NodeId
	ParticipantId
}
