package mocks

import (
	"context"
	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/bots"
)

type MockSenderService struct {
}

func NewMockSenderService() *MockSenderService {
	return &MockSenderService{}
}

func (s *MockSenderService) Send(ctx context.Context, block *bots.Block) error {
	return nil
}
