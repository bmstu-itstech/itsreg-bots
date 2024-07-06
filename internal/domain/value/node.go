package value

import (
	"errors"
)

var (
	ErrInvalidNode     = errors.New("invalid node")
	ErrIncorrectAnswer = errors.New("incorrect answer")
)

type Node struct {
	Type    NodeType
	State   State
	Default State
	Options []Option
}

func NewMessageNode(state State, def State) (Node, error) {
	if state == StateNone {
		return Node{}, ErrInvalidNode
	}

	if def == StateNone {
		return Node{}, ErrInvalidNode
	}

	return Node{
		Type:    Message,
		State:   state,
		Default: def,
		Options: []Option{},
	}, nil
}

func NewQuestionNode(state State, def State) (Node, error) {
	if state == StateNone {
		return Node{}, ErrInvalidNode
	}

	if def == StateNone {
		return Node{}, ErrInvalidNode
	}

	return Node{
		Type:    Question,
		State:   state,
		Default: def,
		Options: []Option{},
	}, nil
}

func NewSelectionNode(state State, options []Option) (Node, error) {
	if state == StateNone {
		return Node{}, ErrInvalidNode
	}

	if options == nil || len(options) == 0 {
		return Node{}, ErrInvalidNode
	}

	return Node{
		Type:    Selection,
		State:   state,
		Default: StateNone,
		Options: options,
	}, nil
}

func (n Node) Next(s string) (State, error) {
	if len(n.Options) > 0 {
		for _, opt := range n.Options {
			if opt.Match(s) {
				return opt.Next, nil
			}
		}
		return n.State, ErrIncorrectAnswer
	}

	return n.Default, nil
}
