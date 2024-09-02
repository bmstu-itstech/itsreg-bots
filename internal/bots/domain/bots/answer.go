package bots

import (
	"errors"
)

type Answer struct {
	State int
	Text  string
}

func NewAnswer(state int, text string) (Answer, error) {
	if state == 0 {
		return Answer{}, errors.New("missing state")
	}

	if text == "" {
		return Answer{}, errors.New("missing text")
	}

	return Answer{
		State: state,
		Text:  text,
	}, nil
}

func MustNewAnswer(state int, text string) Answer {
	a, err := NewAnswer(state, text)
	if err != nil {
		panic(err)
	}
	return a
}
