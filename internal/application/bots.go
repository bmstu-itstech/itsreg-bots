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

func (s *BotsService) Create(
	ctx context.Context,
	name string,
	token string,
	start uint64,
	blocks []dto.Block,
) (uint64, error) {
	return s.manager.Create(ctx, name, token, start, blocks)
}

func (s *BotsService) Token(
	ctx context.Context,
	botId uint64,
) (string, error) {
	return s.manager.Token(ctx, botId)
}
