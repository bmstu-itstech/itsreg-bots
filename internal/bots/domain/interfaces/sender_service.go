package interfaces

import (
	"context"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
)

type SenderService interface {
	Send(
		ctx context.Context,
		botUUID string,
		userID int64,
		msg bots.Message,
	) error
}
