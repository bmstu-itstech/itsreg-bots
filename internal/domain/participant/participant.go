package participant

import (
	"errors"
	"github.com/zhikh23/itsreg-tgservice/internal/domain/module"
)

var (
	ErrInvalidParticipantId = errors.New("invalid participant id")
)

// Participant is any user who uses a bot.
type Participant struct {
	Id      int32          // Unique Participant id
	Current *module.Module // Current Module
}

// New returns Participant and nil, if args are valid; otherwise nil and error.
func New(id int32) (*Participant, error) {
	if id <= 0 {
		return nil, ErrInvalidParticipantId
	}
	return &Participant{
		Id:      id,
		Current: nil,
	}, nil
}
