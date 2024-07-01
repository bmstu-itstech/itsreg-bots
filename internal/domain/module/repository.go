package module

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/entity"
	"github.com/zhikh23/itsreg-bots/internal/objects"
)

var (
	ErrModuleAlreadyExists = errors.New("module already exists")
	ErrModuleNotFound      = errors.New("module not found")
)

type Repository interface {
	Save(module entity.Module) error
	Get(botId int64, nodeId objects.State) (entity.Module, error)
}
