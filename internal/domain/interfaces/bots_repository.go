package interfaces

import (
	"context"
	"fmt"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type BotNotFoundError struct {
	uuid string
}

func (e BotNotFoundError) Error() string {
	return fmt.Sprintf("bot not found: %s", e.uuid)
}

type BotAlreadyExistsError struct {
	uuid string
}

func (e BotAlreadyExistsError) Error() string {
	return fmt.Sprintf("bot already exists: %s", e.uuid)
}

type BotsRepository interface {
	Save(ctx context.Context, bot *bots.Bot) error
	Bot(ctx context.Context, uuid string) (*bots.Bot, error)
	Update(
		ctx context.Context,
		uuid string,
		updateFn func(context.Context, *bots.Bot) (*bots.Bot, error),
	) error
	Delete(ctx context.Context, uuid string) error
}
