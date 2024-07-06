package value

type UserId uint64

const UnknownUserId UserId = 0

func (id UserId) IsUnknown() bool {
	return id == UnknownUserId
}
