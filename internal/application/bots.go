package application

import (
	"context"
	"github.com/zhikh23/itsreg-bots/internal/domain/service"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type BotsService struct {
	processor service.Processor
}

func NewBotsService(processor service.Processor) *BotsService {
	return &BotsService{
		processor: processor,
	}
}

func (s *BotsService) Process(
	ctx context.Context,
	botId uint64,
	userId uint64,
	ans string,
) ([]service.Message, error) {
	return s.processor.Process(ctx, value.BotId(botId), value.UserId(userId), ans)
}
