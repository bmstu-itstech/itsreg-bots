package sender

import (
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/objects"
)

type Sender interface {
	SendMessage(receiver entity.ParticipantId, msg string, buttons []objects.Button) error
}
