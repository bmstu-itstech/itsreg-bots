package bots

import "github.com/bmstu-itstech/itsreg-bots/internal/common/commonerrs"

type Mailing struct {
	Name         string
	EntryKey     string
	RequireState int
}

func NewMailing(
	name string,
	entryKey string,
	requireState int,
) (Mailing, error) {
	if name == "" {
		return Mailing{}, commonerrs.NewInvalidInputError("expected not empty name")
	}

	if entryKey == "" {
		return Mailing{}, commonerrs.NewInvalidInputError("expected not empty entryKey")
	}

	return Mailing{
		Name:         name,
		EntryKey:     entryKey,
		RequireState: requireState,
	}, nil
}

func MustNewMailing(
	name string,
	entryKey string,
	requireState int,
) Mailing {
	m, err := NewMailing(name, entryKey, requireState)
	if err != nil {
		panic(err)
	}
	return m
}

func (m Mailing) IsZero() bool {
	return m == Mailing{}
}
