package bots

import "errors"

type Option struct {
	Text string
	Next int
}

func (o Option) IsZero() bool {
	return o == Option{}
}

func NewOption(text string, next int) (Option, error) {
	if text == "" {
		return Option{}, errors.New("missing text")
	}

	return Option{
		Text: text,
		Next: next,
	}, nil
}

func MustNewOption(text string, next int) Option {
	o, err := NewOption(text, next)
	if err != nil {
		panic(err)
	}
	return o
}

func (o Option) Match(text string) bool {
	return o.Text == text
}
