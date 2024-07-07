package entity

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type Participant struct {
	Id      value.ParticipantId
	Current value.State
}

func NewParticipant(id value.ParticipantId, current value.State) (*Participant, error) {
	return &Participant{
		Id:      id,
		Current: current,
	}, nil
}

func (p *Participant) SwitchTo(node value.Node) {
	p.Current = node.State
}
