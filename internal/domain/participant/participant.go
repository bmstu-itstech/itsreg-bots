package participant

import (
	"errors"
)

var (
	ErrInvalidParticipantId      = errors.New("invalid participant id; should not be zero")
	ErrInvalidParticipantCurrent = errors.New("invalid participant current module id; should not be zero")
)

// Participant is any user who uses a bot.
type Participant struct {
	Id      int32 // Unique Participant id
	Current int32 // Current Module id
}

// New returns Participant and nil, if args are valid; otherwise nil and error.
func New(id int32, current int32) (*Participant, []error) {
	errs := make([]error, 0)
	var err error

	err = validateParticipantId(id)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateParticipantCurrent(current)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Participant{
		Id:      id,
		Current: current,
	}, nil
}

// validateParticipantId return an error, if the participant's id is invalid; otherwise returns nil.
func validateParticipantId(id int32) error {
	if id <= 0 {
		return ErrInvalidParticipantId
	}
	return nil
}

// validateParticipantCurrent return an error, if the participant's current module id is invalid; otherwise returns nil.
func validateParticipantCurrent(current int32) error {
	if current <= 0 {
		return ErrInvalidParticipantCurrent
	}
	return nil
}
