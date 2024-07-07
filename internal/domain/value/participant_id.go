package value

import "errors"

var (
	ErrInvalidParticipantBotId  = errors.New("invalid participant bot id")
	ErrInvalidParticipantUserId = errors.New("invalid participant user id")
)

type ParticipantId struct {
	BotId  BotId
	UserId UserId
}

func NewParticipantId(botId BotId, userId UserId) (ParticipantId, error) {
	if botId.IsUnknown() {
		return ParticipantId{}, ErrInvalidParticipantBotId
	}

	if userId.IsUnknown() {
		return ParticipantId{}, ErrInvalidParticipantUserId
	}

	return ParticipantId{
		BotId:  botId,
		UserId: userId,
	}, nil
}
