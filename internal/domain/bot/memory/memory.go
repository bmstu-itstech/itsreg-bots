package memory

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/bot"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"sync"
)

type Repository struct {
	bots map[int64]entity.Bot
	sync.Mutex
}

func New() *Repository {
	return &Repository{
		bots: make(map[int64]entity.Bot),
	}
}

func (r *Repository) Save(b entity.Bot) error {
	if _, ok := r.bots[b.Id]; ok {
		return bot.ErrBotAlreadyExists
	}

	r.Lock()
	r.bots[b.Id] = b
	r.Unlock()

	return nil
}

func (r *Repository) Get(id int64) (entity.Bot, error) {
	if b, ok := r.bots[id]; ok {
		return b, nil
	}

	return entity.Bot{}, bot.ErrBotNotFound
}
