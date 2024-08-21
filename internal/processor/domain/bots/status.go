package bots

import (
	"fmt"
	cerrors "github.com/bmstu-itstech/itsreg-bots/internal/common/errors"
)

type Status struct {
	s string
}

var (
	Started = Status{s: "started"}
	Stopped = Status{s: "stopped"}
	Failed  = Status{s: "failed"}
)

func (s Status) IsZero() bool {
	return s == Status{}
}

func (s Status) String() string {
	return s.s
}

func NewStatusFromString(s string) (Status, error) {
	switch s {
	case "started":
		return Started, nil
	case "stopped":
		return Stopped, nil
	case "failed":
		return Failed, nil
	}
	return Status{}, cerrors.NewIncorrectInputError(
		fmt.Sprintf("invalid status: %s", s),
		"invalid-status",
	)
}
