package bots

import "context"

type RunnerPublisher interface {
	PublishStart(ctx context.Context, botUUID string) error
	PublishStop(ctx context.Context, botUUID string) error
}
