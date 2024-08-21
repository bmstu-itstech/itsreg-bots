package bots

import "errors"

type Participant struct {
	UserID int64
	State  int
}

func NewParticipant(
	id int64,
	state int,
) (*Participant, error) {
	if id == 0 {
		return nil, errors.New("missing id")
	}

	return &Participant{
		UserID: id,
		State:  state,
	}, nil
}

func (p *Participant) IsFinished() bool {
	return p.State == FinishState
}

func (p *Participant) SwitchTo(state int) {
	p.State = state
}

func MustNewParticipant(
	id int64,
	state int,
) *Participant {
	p, err := NewParticipant(id, state)
	if err != nil {
		panic(err)
	}
	return p
}
