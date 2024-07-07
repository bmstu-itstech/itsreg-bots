package value

import (
	"errors"
	"github.com/gabrielmellooliveira/go-spec"
)

var (
	ErrInvalidMessageNode   = errors.New("invalid message node")
	ErrInvalidQuestionNode  = errors.New("invalid question node")
	ErrInvalidSelectionNode = errors.New("invalid selection node")
	ErrIncorrectAnswer      = errors.New("incorrect answer")
)

type Node struct {
	Type    NodeType
	State   State
	Default State
	Options []Option
}

func NewNode(
	typ NodeType,
	state State,
	def State,
	options []Option,
) (Node, error) {
	node := Node{
		Type:    typ,
		State:   state,
		Default: def,
		Options: options,
	}

	switch typ {
	case Message:
		s := newMessageNodeSpec()
		if !s.IsSatisfiedBy(node) {
			return Node{}, ErrInvalidMessageNode
		}
	case Question:
		s := newQuestionNodeSpec()
		if !s.IsSatisfiedBy(node) {
			return Node{}, ErrInvalidQuestionNode
		}
	case Selection:
		s := newSelectionNodeSpec()
		if !s.IsSatisfiedBy(node) {
			return Node{}, ErrInvalidSelectionNode
		}
	default:
		return Node{}, ErrInvalidNodeType
	}

	return node, nil
}

func NewMessageNode(state State, def State) (Node, error) {
	node := Node{
		Type:    Message,
		State:   state,
		Default: def,
		Options: make([]Option, 0),
	}

	s := newMessageNodeSpec()
	if !s.IsSatisfiedBy(node) {
		return Node{}, ErrInvalidMessageNode
	}

	return node, nil
}

func NewQuestionNode(state State, def State) (Node, error) {
	node := Node{
		Type:    Question,
		State:   state,
		Default: def,
		Options: make([]Option, 0),
	}

	s := newQuestionNodeSpec()
	if !s.IsSatisfiedBy(node) {
		return Node{}, ErrInvalidQuestionNode
	}

	return node, nil
}

func NewSelectionNode(state State, options []Option) (Node, error) {
	node := Node{
		Type:    Selection,
		State:   state,
		Default: StateNone,
		Options: options,
	}

	s := newSelectionNodeSpec()
	if !s.IsSatisfiedBy(node) {
		return Node{}, ErrInvalidSelectionNode
	}

	return node, nil
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

type nodeSpec spec.Specification[Node]

type nodeHasValidStateSpec struct {
	spec.BaseSpecification[Node]
}

func newNodeHasValidStateSpec() *nodeHasValidStateSpec {
	return &nodeHasValidStateSpec{}
}

func (s nodeHasValidStateSpec) IsSatisfiedBy(node Node) bool {
	return !node.State.IsNone()
}

type nodeHasDefaultSpec struct {
	spec.BaseSpecification[Node]
}

func newNodeHasDefaultSpec() *nodeHasDefaultSpec {
	return &nodeHasDefaultSpec{}
}

func (s nodeHasDefaultSpec) IsSatisfiedBy(node Node) bool {
	return !node.Default.IsNone()
}

type nodeHasOptionsSpec struct {
	spec.BaseSpecification[Node]
}

func newNodeHasOptionsSpec() *nodeHasOptionsSpec {
	return &nodeHasOptionsSpec{}
}

func (s nodeHasOptionsSpec) IsSatisfiedBy(node Node) bool {
	return len(node.Options) > 0
}

func newMessageNodeSpec() nodeSpec {
	return newNodeHasValidStateSpec().
		And(newNodeHasDefaultSpec()).
		Not(newNodeHasOptionsSpec())
}

func newQuestionNodeSpec() nodeSpec {
	return newNodeHasValidStateSpec().
		And(newNodeHasDefaultSpec()).
		Not(newNodeHasOptionsSpec())
}

func newSelectionNodeSpec() nodeSpec {
	return newNodeHasValidStateSpec().
		And(newNodeHasOptionsSpec()).
		Not(newNodeHasDefaultSpec())
}
