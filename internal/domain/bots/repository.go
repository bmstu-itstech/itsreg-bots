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

type Repository interface {
	UpdateOrCreate(ctx context.Context, bot *Bot) error
	UpdateStatus(ctx context.Context, botUUID string, status Status) error
	Delete(ctx context.Context, uuid string) error

	Bot(ctx context.Context, uuid string) (*Bot, error)
	UserBots(ctx context.Context, userUUID string) ([]*Bot, error)
}
