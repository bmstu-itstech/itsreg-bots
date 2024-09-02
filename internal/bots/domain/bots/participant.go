package bots

import (
	"github.com/bmstu-itstech/itsreg-bots/internal/common/commonerrs"
)

type Participant struct {
	BotUUID string
	UserID  int64
	State   int
	answers map[int]Answer
}

func NewParticipant(
	botUUID string,
	id int64,
) (*Participant, error) {
	if botUUID == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty botUUID")
	}

	if id == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty id")
	}

	return &Participant{
		BotUUID: botUUID,
		UserID:  id,
		State:   0,
		answers: make(map[int]Answer),
	}, nil
}

func MustNewParticipant(
	botUUID string,
	id int64,
) *Participant {
	p, err := NewParticipant(botUUID, id)
	if err != nil {
		panic(err)
	}
	return p
}

func NewParticipantFromDB(
	botUUID string,
	id int64,
	state int,
	answers []Answer,
) (*Participant, error) {
	if botUUID == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty botUUID")
	}

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
		answers: m,
	}, nil
}

func (p *Participant) Answers() []Answer {
	answers := make([]Answer, 0, len(p.answers))
	for _, a := range p.answers {
		answers = append(answers, a)
	}
	return answers
}

func (p *Participant) IsProcessing() bool {
	return p.State != 0
}

func (p *Participant) SwitchTo(state int) {
	p.State = state
}

func (p *Participant) AddAnswer(text string) error {
	ans, err := NewAnswer(p.UserID, p.State, text)
	if err != nil {
		return err
	}
	p.answers[p.State] = ans
	return nil
}

func (p *Participant) CleanAnswerIfExists(state int) {
	if _, ok := p.answers[state]; ok {
		delete(p.answers, state)
	}
}
