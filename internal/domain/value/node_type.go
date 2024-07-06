package value

import "errors"

type NodeType int

var (
	ErrInvalidNodeType = errors.New("invalid node type")
)

const (
	Message NodeType = iota + 1
	Question
	Selection
)

func NewNodeType(i int) (NodeType, error) {
	t := NodeType(i)
	switch t {
	case Message, Question, Selection:
		return t, nil
	}
	return 0, ErrInvalidNodeType
}
