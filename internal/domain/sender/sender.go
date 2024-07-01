package sender

import "github.com/zhikh23/itsreg-bots/internal/entity"

type Sender interface {
	SendMessage(receiver entity.ParticipantId, msg string, buttons []entity.Button) error
}
