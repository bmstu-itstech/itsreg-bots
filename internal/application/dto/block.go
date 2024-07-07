package dto

import (
	"github.com/zhikh23/itsreg-bots/internal/domain/entity"
	"github.com/zhikh23/itsreg-bots/internal/domain/value"
)

type Block struct {
	State   uint64
	Type    int
	Default uint64
	Options []Option
	Title   string
	Text    string
}

func BlockFromDto(dto Block) (*entity.Block, error) {
	typ, err := value.NewNodeType(dto.Type)
	if err != nil {
		return nil, err
	}

	var node value.Node

	switch typ {
	case value.Message:
		node, err = value.NewMessageNode(value.State(dto.State), value.State(dto.Default))
		if err != nil {
			return nil, err
		}
		if len(dto.Options) > 0 {
			return nil, value.ErrInvalidMessageNode
		}
	case value.Question:
		node, err = value.NewQuestionNode(value.State(dto.State), value.State(dto.Default))
		if err != nil {
			return nil, err
		}
		if len(dto.Options) > 0 {
			return nil, value.ErrInvalidQuestionNode
		}
	case value.Selection:
		options, err := OptionsFromDto(dto.Options)
		if err != nil {
			return nil, err
		}
		node, err = value.NewSelectionNode(value.State(dto.State), options)
		if err != nil {
			return nil, err
		}
	}

	block, err := entity.NewBlock(node, value.UnknownBotId, dto.Title, dto.Text)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func BlocksFromDtos(dtos []Block) ([]*entity.Block, error) {
	blocks := make([]*entity.Block, len(dtos))
	for i, dto := range dtos {
		block, err := BlockFromDto(dto)
		if err != nil {
			return nil, err
		}
		blocks[i] = block
	}
	return blocks, nil
}
