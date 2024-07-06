package entity

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

var (
	ErrInvalidParticipant = errors.New("invalid participant")
)

type Participant struct {
	Id      value.ParticipantId
	Current value.State
}

func NewParticipant(id value.ParticipantId, current value.State) (*Participant, error) {
	if id.BotId == 0 {
		return nil, ErrInvalidParticipant
	}

	if id.UserId == 0 {
		return nil, ErrInvalidParticipant
	}

	if current.IsNone() {
		return nil, ErrInvalidParticipant
	}

	return &Participant{
		Id:      id,
		Current: current,
	}, nil
}

func (p *Participant) SwitchTo(node value.Node) {
	p.Current = node.State
}
