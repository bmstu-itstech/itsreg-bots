package entity

type ParticipantId struct {
	BotId  int64
	UserId int64
}

type Participant struct {
	CurrentId NodeId
	ParticipantId
}
