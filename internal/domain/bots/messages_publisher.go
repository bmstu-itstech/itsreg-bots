package bots

import (
	"context"
)

type MessagesPublisher interface {
	Publish(
		ctx context.Context,
		botUUID string,
		userID int64,
		msg Message,
	) error
}
