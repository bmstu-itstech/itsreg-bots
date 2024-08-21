package bots

import (
	"errors"
	"fmt"
	"time"

	cerrors "github.com/bmstu-itstech/itsreg-bots/internal/common/errors"
)

type Bot struct {
	UUID string

	Blocks     map[int]Block
	StartState int

	Name  string
	Token string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBot(
	uuid string,
	blocks []Block,
	startState int,
	name string,
	token string,
) (*Bot, error) {
	if uuid == "" {
		return nil, errors.New("missing uuid")
	}

	if blocks == nil {
		return nil, errors.New("missing blocks")
	}

	for _, block := range blocks {
		if block.IsZero() {
			return nil, fmt.Errorf("block with state %d is empty", block.State)
		}
	}

	if startState == 0 {
		return nil, errors.New("missing start state")
	}

	if name == "" {
		return nil, errors.New("missing name")
	}

	if token == "" {
		return nil, errors.New("missing token")
	}

	mappedBlocks, err := mapBlocks(blocks)
	if err != nil {
		return nil, err
	}

	//err = traverse(startState, mappedBlocks)
	//if err != nil {
	//	return nil, err
	//}

	return &Bot{
		UUID:       uuid,
		Blocks:     mappedBlocks,
		StartState: startState,
		Name:       name,
		Token:      token,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func mapBlocks(blocks []Block) (map[int]Block, error) {
	mapped := map[int]Block{}
	for _, block := range blocks {
		if _, ok := mapped[block.State]; ok {
			return nil, cerrors.NewIncorrectInputError(
				fmt.Sprintf("block with state %d is dublicated", block.State),
				"invalid-blocks",
			)
		}
		mapped[block.State] = block
	}
	return mapped, nil
}

type color int

const (
	white color = iota
	grey
	black
)

type vertex struct {
	Block Block
	Color color
}

func traverse(
	startState int,
	blocks map[int]Block,
) error {
	vertices := map[int]vertex{}

	for state, block := range blocks {
		vertices[state] = vertex{
			Block: block,
			Color: white,
		}
	}

	start, ok := vertices[startState]
	if !ok {
		return cerrors.NewIncorrectInputError(
			fmt.Sprintf("block with state %d not found", startState),
			"invalid-blocks",
		)
	}

	err := dfs(start, vertices)
	if err != nil {
		return err
	}

	for state, v := range vertices {
		if v.Color == white {
			return cerrors.NewIncorrectInputError(
				fmt.Sprintf("block with state %d is unused", state),
				"invalid-blocks",
			)
		}
	}

	return nil
}

func dfs(
	current vertex,
	vertices map[int]vertex,
) error {
	current.Color = grey

	for _, opt := range current.Block.Options {
		nextState := opt.Next
		if nextState == 0 {
			continue
		}

		next, ok := vertices[nextState]
		if !ok {
			return cerrors.NewIncorrectInputError(
				fmt.Sprintf("block with state %d not found", nextState),
				"invalid-blocks",
			)
		}

		if next.Color == white {
			err := dfs(next, vertices)
			if err != nil {
				return err
			}
		} else {
			return cerrors.NewIncorrectInputError(
				"recursion detected",
				"invalid-blocks",
			)
		}
		next.Color = black
	}

	return nil
}

func MustNewBot(
	uuid string,
	blocks []Block,
	startState int,
	name string,
	token string,
) *Bot {
	b, err := NewBot(uuid, blocks, startState, name, token)
	if err != nil {
		panic(err)
	}
	return b
}

func NewBotFromDB(
	uuid string,
	blocks []Block,
	startState int,
	name string,
	token string,
	createdAt time.Time,
	updatedAt time.Time,
) (*Bot, error) {
	if uuid == "" {
		return nil, errors.New("missing uuid")
	}

	if blocks == nil {
		return nil, errors.New("missing blocks")
	}

	mappedBlocks := map[int]Block{}
	for _, block := range blocks {
		mappedBlocks[block.State] = block
	}

	if startState == 0 {
		return nil, errors.New("missing start state")
	}

	if name == "" {
		return nil, errors.New("missing name")
	}

	if token == "" {
		return nil, errors.New("missing token")
	}

	if createdAt.IsZero() {
		return nil, errors.New("missing created at")
	}

	if updatedAt.IsZero() {
		return nil, errors.New("missing updated at")
	}

	return &Bot{
		UUID:       uuid,
		Blocks:     mappedBlocks,
		StartState: startState,
		Name:       name,
		Token:      token,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
