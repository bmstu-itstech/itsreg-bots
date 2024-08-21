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
		Type:      QuestionBlock,
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
