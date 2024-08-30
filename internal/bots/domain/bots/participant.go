package bots

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/common/commonerrs"
)

type Participant struct {
	UserID  int64
	State   int
	Answers map[int]Answer
}

func NewParticipant(
	id int64,
	state int,
) (*Participant, error) {
	if id == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty id")
	}

	return &Participant{
		UserID:  id,
		State:   state,
		Answers: make(map[int]Answer),
	}, nil
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

func NewParticipantFromDB(
	id int64,
	state int,
	answers []Answer,
) (*Participant, error) {
	if id == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty id")
	}

	if answers == nil {
		answers = make([]Answer, 0)
	}

	m := make(map[int]Answer)
	for _, a := range answers {
		m[a.State] = a
	}

	return &Participant{
		UserID:  id,
		State:   state,
		Answers: m,
	}, nil
}

func (p *Participant) IsFinished() bool {
	return p.State == FinishState
}

func (p *Participant) SwitchTo(state int) {
	p.State = state
}

func (p *Participant) AddAnswer(a Answer) {
	p.Answers[a.State] = a
}
