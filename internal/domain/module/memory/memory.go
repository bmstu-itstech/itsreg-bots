package memory

import (
	module2 "github.com/zhikh23/itsreg-bots/internal/domain/module"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/objects"
	"sync"
)

type pair struct {
	botId  int64
	nodeId objects.State
}

type Repository struct {
	modules map[pair]entity.Module
	sync.Mutex
}

func New() *Repository {
	return &Repository{
		modules: make(map[pair]entity.Module),
	}
}

func (r *Repository) Save(module entity.Module) error {
	key := pair{
		botId:  module.BotId,
		nodeId: module.Id,
	}

	if _, ok := r.modules[key]; ok {
		return module2.ErrModuleAlreadyExists
	}

	r.Lock()
	r.modules[key] = module
	r.Unlock()

	return nil
}

func (r *Repository) Get(botId int64, nodeId objects.State) (entity.Module, error) {
	key := pair{
		botId:  botId,
		nodeId: nodeId,
	}

	if module, ok := r.modules[key]; ok {
		return module, nil
	}

	return entity.Module{}, module2.ErrModuleNotFound
}
