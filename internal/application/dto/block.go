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

func BlockFromDto(dto Block, botId value.BotId) (*entity.Block, error) {
	typ, err := value.NewNodeType(dto.Type)
	if err != nil {
		return nil, err
	}

	options, err := OptionsFromDto(dto.Options)
	if err != nil {
		return nil, err
	}

	node, err := value.NewNode(typ, value.State(dto.State), value.State(dto.Default), options)
	if err != nil {
		return nil, err
	}

	block, err := entity.NewBlock(node, botId, dto.Title, dto.Text)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func BlocksFromDtos(dtos []Block, botId value.BotId) ([]*entity.Block, error) {
	blocks := make([]*entity.Block, len(dtos))
	for i, dto := range dtos {
		block, err := BlockFromDto(dto, botId)
		if err != nil {
			return nil, err
		}
		blocks[i] = block
	}
	return blocks, nil
}
