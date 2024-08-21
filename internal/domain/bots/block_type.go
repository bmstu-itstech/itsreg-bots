package bots

import (
	"fmt"

	"github.com/bmstu-itstech/itsreg-bots/internal/common/errors"
)

type BlockType struct {
	s string
}

var (
	MessageBlock   = BlockType{s: "message"}
	QuestionBlock  = BlockType{s: "question"}
	SelectionBlock = BlockType{s: "selection"}
)

func (b BlockType) String() string {
	return b.s
}

func (b BlockType) IsZero() bool {
	return b == BlockType{}
}

func NewBlockTypeFromString(s string) (BlockType, error) {
	switch s {
	case "message":
		return MessageBlock, nil
	case "question":
		return QuestionBlock, nil
	case "selection":
		return SelectionBlock, nil
	}
	return BlockType{}, errors.NewIncorrectInputError(
		fmt.Sprintf("invalid block type: %s", s),
		"invalid-block-type",
	)
}
