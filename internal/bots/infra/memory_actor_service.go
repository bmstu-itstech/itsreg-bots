package infra

import (
	"context"

	"github.com/bmstu-itstech/itsreg-bots/internal/bots/domain/interfaces"
)

type memoryActorService struct {
	startChan chan struct{}
	stopChan  chan struct{}
}

func NewMemoryActorService(startChan, stopChan chan struct{}) interfaces.ActorService {
	return memoryActorService{
		startChan: startChan,
		stopChan:  stopChan,
	}
}

func (a memoryActorService) Start(ctx context.Context, botUUID string) error {
	a.startChan <- struct{}{}
	return nil
}

func (a memoryActorService) Stop(ctx context.Context, botUUID string) error {
	a.stopChan <- struct{}{}
	return nil
}
