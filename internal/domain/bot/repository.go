package bot

import (
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/entity"
)

var (
	ErrBotNotFound      = errors.New("bot not found")
	ErrBotAlreadyExists = errors.New("bot already exists")
)

type Repository interface {
	Save(b entity.Bot) error
	Get(id int64) (entity.Bot, error)
}
