package value

type BotId uint64

const UnknownBotId BotId = 0

func (id BotId) IsUnknown() bool {
	return id == UnknownBotId
}
