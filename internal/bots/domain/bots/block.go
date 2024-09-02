package bots

import "errors"

type Block struct {
	Type      BlockType
	State     int
	NextState int
	Options   []Option

	Title string
	Text  string
}

func (b Block) IsZero() bool {
	return b.Type.IsZero()
}

func NewBlock(
	blockType string,
	state int,
	nextState int,
	options []Option,
	title string,
	text string,
) (Block, error) {
	bt, err := NewBlockTypeFromString(blockType)
	if err != nil {
		return Block{}, err
	}

	switch bt {
	case MessageBlock:
		return NewMessageBlock(state, nextState, title, text)
	case QuestionBlock:
		return NewQuestionBlock(state, nextState, title, text)
	case SelectionBlock:
		return NewSelectionBlock(state, nextState, options, title, text)
	}
	return Block{}, errors.New("unknown type")
}

func NewMessageBlock(
	state int,
	next int,
	title string,
	text string,
) (Block, error) {
	if state == 0 {
		return Block{}, errors.New("missing state")
	}

	if title == "" {
		return Block{}, errors.New("missing title")
	}

	if text == "" {
		return Block{}, errors.New("missing text")
	}

	return Block{
		Type:      MessageBlock,
		State:     state,
		NextState: next,
		Options:   nil,
		Title:     title,
		Text:      text,
	}, nil
}

func MustNewMessageBlock(
	state int,
	next int,
	title string,
	text string,
) Block {
	b, err := NewMessageBlock(state, next, title, text)
	if err != nil {
		panic(err)
	}
	return b
}

func NewQuestionBlock(
	state int,
	next int,
	title string,
	text string,
) (Block, error) {
	if state == 0 {
		return Block{}, errors.New("missing state")
	}

	if title == "" {
		return Block{}, errors.New("missing title")
	}

	if text == "" {
		return Block{}, errors.New("missing text")
	}

	return Block{
		Type:      QuestionBlock,
		State:     state,
		NextState: next,
		Options:   nil,
		Title:     title,
		Text:      text,
	}, nil
}

func MustNewQuestionBlock(
	state int,
	next int,
	title string,
	text string,
) Block {
	b, err := NewQuestionBlock(state, next, title, text)
	if err != nil {
		panic(err)
	}
	return b
}

func NewSelectionBlock(
	state int,
	next int,
	options []Option,
	title string,
	text string,
) (Block, error) {
	if state == 0 {
		return Block{}, errors.New("missing state")
	}

	if len(options) == 0 {
		return Block{}, errors.New("missing options")
	}

	if title == "" {
		return Block{}, errors.New("missing title")
	}

	if text == "" {
		return Block{}, errors.New("missing text")
	}

	return Block{
		Type:      SelectionBlock,
		State:     state,
		NextState: next,
		Options:   options,
		Title:     title,
		Text:      text,
	}, nil
}

func MustNewSelectionBlock(
	state int,
	next int,
	options []Option,
	title string,
	text string,
) Block {
	b, err := NewSelectionBlock(state, next, options, title, text)
	if err != nil {
		panic(err)
	}
	return b
}

func NewBlockFromDB(
	blockType string,
	state int,
	next int,
	options []Option,
	title string,
	text string,
) (Block, error) {
	if blockType == "" {
		return Block{}, errors.New("missing block type")
	}

	t, err := NewBlockTypeFromString(blockType)
	if err != nil {
		return Block{}, err
	}

	if state == 0 {
		return Block{}, errors.New("missing state")
	}

	if title == "" {
		return Block{}, errors.New("missing title")
	}

	if text == "" {
		return Block{}, errors.New("missing text")
	}

	return Block{
		Type:      t,
		State:     state,
		NextState: next,
		Options:   options,
		Title:     title,
		Text:      text,
	}, nil
}

func (b Block) Message() (Message, error) {
	switch b.Type {
	case MessageBlock, QuestionBlock:
		return NewPlainMessage(b.Text)
	case SelectionBlock:
		return NewMessageWithButtons(b.Text, b.Options)
	}
	return Message{}, errors.New("unknown type")
}

func (b Block) IsFinish() bool {
	return len(b.Options) == 0 && b.NextState == 0
}

func (b Block) ChildrenStates() []int {
	states := make([]int, 0)
	if b.NextState != 0 {
		states = append(states, b.NextState)
	}

	for _, opt := range b.Options {
		states = append(states, opt.Next)
	}

	return states
}
