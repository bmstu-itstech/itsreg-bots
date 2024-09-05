package bots

import "github.com/bmstu-itstech/itsreg-bots/internal/common/commonerrs"

type EntryPoint struct {
	Key   string
	State int
}

func (e EntryPoint) IsZero() bool {
	return e == EntryPoint{}
}

func NewEntryPoint(key string, state int) (EntryPoint, error) {
	if key == "" {
		return EntryPoint{}, commonerrs.NewInvalidInputError("expected not empty entrypoint key")
	}

	if state == 0 {
		return EntryPoint{}, commonerrs.NewInvalidInputError("expected not empty entrypoint state")
	}

	return EntryPoint{
		Key:   key,
		State: state,
	}, nil
}

func MustNewEntryPoint(key string, state int) EntryPoint {
	e, err := NewEntryPoint(key, state)
	if err != nil {
		panic(err)
	}
	return e
}
