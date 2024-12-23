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

func (r *mockBotRepository) Update(
	ctx context.Context,
	botUUID string,
	updateFn func(innerCtx context.Context, bot *bots.Bot) error,
) error {
	r.Lock()
	defer r.Unlock()

	bot, ok := r.m[botUUID]
	if !ok {
		return bots.BotNotFoundError{UUID: botUUID}
	}

	err := updateFn(ctx, &bot)
	if err != nil {
		return err
	}

	r.m[botUUID] = bot

	return nil
}

func (r *mockBotRepository) UpdateOrCreate(_ context.Context, bot *bots.Bot) error {
	r.Lock()
	defer r.Unlock()

	r.m[bot.UUID] = *bot

	return nil
}

func (r *mockBotRepository) UpdateStatus(_ context.Context, uuid string, status bots.Status) error {
	r.Lock()
	defer r.Unlock()

	bot, ok := r.m[uuid]
	if !ok {
		return bots.BotNotFoundError{UUID: uuid}
	}

	bot.SetStatus(status)
	r.m[uuid] = bot

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

func (r *mockBotRepository) Bot(_ context.Context, uuid string) (*bots.Bot, error) {
	r.RLock()
	defer r.RUnlock()

	b, ok := r.m[uuid]
	if !ok {
		return nil, bots.BotNotFoundError{UUID: uuid}
	}

	return &b, nil
}

func (r *mockBotRepository) UserBots(_ context.Context, userUUID string) ([]*bots.Bot, error) {
	r.RLock()
	defer r.RUnlock()

	return r.botsFilter(func(bot *bots.Bot) bool {
		return bot.OwnerUUID == userUUID
	}), nil

}

func (r *mockBotRepository) BotsWithStatus(_ context.Context, status bots.Status) ([]*bots.Bot, error) {
	r.RLock()
	defer r.RUnlock()

	return r.botsFilter(func(bot *bots.Bot) bool {
		return bot.Status == status
	}), nil
}

type botPredicate func(bot *bots.Bot) bool

func (r *mockBotRepository) botsFilter(predicate botPredicate) []*bots.Bot {
	res := make([]*bots.Bot, 0)
	for _, b := range r.m {
		if predicate(&b) {
			res = append(res, &b)
		}
	}
	return res
}
