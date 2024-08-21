package bots

import (
	"errors"
)

type Answer struct {
	UserID int64
	State  int
	Text   string
}

func NewAnswer(userID int64, state int, text string) (*Answer, error) {
	if userID == 0 {
		return nil, errors.New("missing id")
	}

	if state == 0 {
		return nil, errors.New("missing state")
	}

	if text == "" {
		return nil, errors.New("missing text")
	}

	return &Answer{
		UserID: userID,
		State:  state,
		Text:   text,
	}, nil
}
