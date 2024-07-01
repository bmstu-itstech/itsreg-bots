package module

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/entity"
)

var (
	ErrModuleAlreadyExists = errors.New("module already exists")
	ErrModuleNotFound      = errors.New("module not found")
)

type Repository interface {
	Save(module entity.Module) error
	Get(botId int64, nodeId entity.NodeId) (entity.Module, error)
}
