package mocks

import (
	"context"
	"sync"

	"github.com/bmstu-itstech/itsreg-bots/internal/domain/bots"
)

type mockBotRepository struct {
	sync.RWMutex
	m map[string]bots.Bot
}

func NewMockBotRepository() bots.Repository {
	return &mockBotRepository{m: make(map[string]bots.Bot)}
}

func (r *mockBotRepository) Save(_ context.Context, bot *bots.Bot) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[bot.UUID]; ok {
		return bots.BotAlreadyExistsError{UUID: bot.UUID}
	}

	r.m[bot.UUID] = *bot

	return nil
}

func (r *mockBotRepository) Bot(_ context.Context, uuid string) (*bots.Bot, error) {
	r.RLock()
	defer r.RUnlock()

	b, ok := r.m[uuid]
	if !ok {
		return nil, bots.BotNotFoundError{UUID: uuid}
	}

	return &b, nil
}

func (r *mockBotRepository) Bots(_ context.Context, userUUID string) ([]*bots.Bot, error) {
	r.RLock()
	defer r.RUnlock()

	res := make([]*bots.Bot, 0)
	for _, b := range r.m {
		if b.OwnerUUID == userUUID {
			res = append(res, &b)
		}
	}

	return res, nil
}

func (r *mockBotRepository) Update(
	ctx context.Context,
	uuid string,
	updateFn func(context.Context, *bots.Bot) error,
) error {
	r.Lock()
	defer r.Unlock()

	b, ok := r.m[uuid]
	if !ok {
		return bots.BotNotFoundError{UUID: uuid}
	}

	err := updateFn(ctx, &b)
	if err != nil {
		return err
	}

	r.m[uuid] = b

	return nil
}

func (r *mockBotRepository) Delete(_ context.Context, uuid string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[uuid]; !ok {
		return bots.BotNotFoundError{UUID: uuid}
	}

	delete(r.m, uuid)

	return nil
}
