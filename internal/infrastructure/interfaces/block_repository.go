package interfaces

import (
	"context"
	"errors"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/entity"
	"github.com/bmstu-itstech/itsreg-bots/internal/domain/value"
)

var (
	ErrBlockAlreadyExists = errors.New("block already exists")
	ErrBlockNotFound      = errors.New("block not found")
)

type BlockRepository interface {
	Save(context.Context, *entity.Block) error
	Block(context.Context, value.BotId, value.State) (*entity.Block, error)
}
