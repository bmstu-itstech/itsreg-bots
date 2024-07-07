package value

import "errors"

var (
	ErrInvalidOptionText = errors.New("invalid option text")
	ErrInvalidOptionNext = errors.New("invalid option next state")
)

type Option struct {
	Text string
	Next State
}

func NewOption(text string, next State) (Option, error) {
	if len(text) == 0 {
		return Option{}, ErrInvalidOptionText
	}

	if next.IsNone() {
		return Option{}, ErrInvalidOptionNext
	}

	return Option{
		Text: text,
		Next: next,
	}, nil
}

func (o Option) Match(s string) bool {
	return s == o.Text
}
