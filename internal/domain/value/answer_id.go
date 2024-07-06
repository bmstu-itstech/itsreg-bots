package value

import "errors"

var (
	ErrInvalidAnswerId = errors.New("invalid answer id")
)

type AnswerId struct {
	ParticipantId ParticipantId
	State         State
}

func NewAnswerId(participantId ParticipantId, state State) (AnswerId, error) {
	if state.IsNone() {
		return AnswerId{}, ErrInvalidAnswerId
	}

	return AnswerId{
		ParticipantId: participantId,
		State:         state,
	}, nil
}
