package application

import (
	"context"
	"errors"
	"github.com/zhikh23/itsreg-bots/internal/application/dto"
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/interfaces"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
	"log/slog"
)

var (
	ErrRecursionFound   = errors.New("recursion found")
	ErrUnusedBlockFound = errors.New("unused block found")
	ErrNodeNotFound     = errors.New("node not found")
)

type BotsManager struct {
	log     *slog.Logger
	botRepo interfaces.BotRepository
	blcRepo interfaces.BlockRepository
}

func NewBotsManager(
	logger *slog.Logger,
	botRepo interfaces.BotRepository,
	blcRepo interfaces.BlockRepository,
) *BotsManager {
	return &BotsManager{
		log:     logger,
		botRepo: botRepo,
		blcRepo: blcRepo,
	}
}

func (m *BotsManager) Create(
	ctx context.Context,
	name string,
	token string,
	start uint64,
	blocks []dto.Block,
) (uint64, error) {
	const op = "BotsManager.Create"

	log := m.log.With(
		slog.String("op", op))

	log.Info("processing create bot")

	bot, err := entity.NewBot(value.UnknownBotId, name, token, value.State(start))
	if err != nil {
		log.Error("invalid bot's info", "err", err.Error())
		return 0, err
	}

	botId, err := m.botRepo.Save(ctx, bot)
	if err != nil {
		log.Error("failed to save bot", "err", err.Error())
		return 0, err
	}

	eBlocks, err := dto.BlocksFromDtos(blocks, botId)
	if err != nil {
		log.Error("failed to convert blocks from dto", "err", err.Error())
		return 0, err
	}

	err = traverse(value.State(start), eBlocks)
	if err != nil {
		log.Error("invalid bot's script", "err", err.Error())
		return 0, err
	}

	for _, block := range eBlocks {
		err = m.blcRepo.Save(ctx, block)
		if err != nil {
			log.Error("failed to save block", "err", err.Error())
			return 0, err
		}
	}

	return uint64(botId), nil
}

type Color int

const (
	White Color = iota
	Grey
	Black
)

func traverse(
	start value.State,
	blocks []*entity.Block,
) error {
	mBlocks := map[value.State]*entity.Block{}
	colors := map[value.State]Color{}

	for _, block := range blocks {
		mBlocks[block.Node.State] = block
		colors[block.Node.State] = White
	}

	err := dfs(start, mBlocks, colors)
	if err != nil {
		return err
	}

	for _, color := range colors {
		if color == White {
			return ErrUnusedBlockFound
		}
	}

	return nil
}

func dfs(
	current value.State,
	vertices map[value.State]*entity.Block,
	colors map[value.State]Color,
) error {
	colors[current] = Grey

	vertex, ok := vertices[current]
	if !ok {
		return ErrNodeNotFound
	}

	switch vertex.Node.Type {
	case value.Selection:
		for _, opt := range vertex.Node.Options {
			next := opt.Next
			if next.IsNone() {
				continue
			}
			if colors[next] == White {
				err := dfs(next, vertices, colors)
				if err != nil {
					return err
				}
			} else {
				return ErrRecursionFound
			}
			colors[next] = Black
		}

	case value.Message, value.Question:
		next := vertex.Node.Default
		if !next.IsNone() {
			if colors[next] == White {
				err := dfs(next, vertices, colors)
				if err != nil {
					return err
				}
			} else {
				return ErrRecursionFound
			}
			colors[next] = Black
		}
	}
	return nil
}
