package bots

import (
	"context"
	"fmt"
)

type BotNotFoundError struct {
	UUID string
}

func (e BotNotFoundError) Error() string {
	return fmt.Sprintf("bot not found: %s", e.UUID)
}

type BotAlreadyExistsError struct {
	UUID string
}

func (e BotAlreadyExistsError) Error() string {
	return fmt.Sprintf("bot already exists: %s", e.UUID)
}

type Repository interface {
	Save(ctx context.Context, bot *Bot) error
	Bot(ctx context.Context, uuid string) (*Bot, error)
	Bots(ctx context.Context, userUUID string) ([]*Bot, error)
	Update(
		ctx context.Context,
		uuid string,
		updateFn func(context.Context, *Bot) error,
	) error
	Delete(ctx context.Context, uuid string) error
}
