package value

import "errors"

var (
	ErrInvalidOption = errors.New("invalid option")
)

type Option struct {
	Text string
	Next State
}

func NewOption(text string, next State) (Option, error) {
	if len(text) == 0 {
		return Option{}, ErrInvalidOption
	}

	if next.IsNone() {
		return Option{}, ErrInvalidOption
	}

	return Option{
		Text: text,
		Next: next,
	}, nil
}

func (o Option) Match(s string) bool {
	return s == o.Text
}
