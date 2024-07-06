package value

import "errors"

var (
	ErrInvalidParticipantId = errors.New("invalid participant ID")
)

type ParticipantId struct {
	BotId  BotId
	UserId UserId
}

func NewParticipantId(botId BotId, userId UserId) (ParticipantId, error) {
	if botId.IsUnknown() {
		return ParticipantId{}, ErrInvalidParticipantId
	}

	if userId.IsUnknown() {
		return ParticipantId{}, ErrInvalidParticipantId
	}

	return ParticipantId{
		BotId:  botId,
		UserId: userId,
	}, nil
}
