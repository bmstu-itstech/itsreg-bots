package application

import (
	"context"
	"github.com/zhikh23/itsreg-bots/internal/application/dto"
)

type BotsService struct {
	processor *BotsProcessor
	manager   *BotsManager
}

func NewBotsService(
	processor *BotsProcessor,
	manager *BotsManager,
) *BotsService {
	return &BotsService{
		processor: processor,
		manager:   manager,
	}
}

func (s *BotsService) Process(
	ctx context.Context,
	botId uint64,
	userId uint64,
	ans string,
) ([]dto.Message, error) {
	return s.processor.Process(ctx, botId, userId, ans)
}
