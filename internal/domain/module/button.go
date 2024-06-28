package module

import (
	"errors"
)

var (
	ErrEmptyButtonText         = errors.New("invalid button text; should be not empty")
	ErrEmptyButtonValue        = errors.New("invalid button value; should be not empty")
	ErrInvalidButtonNextModule = errors.New("invalid button next module; should be not nil")
)

// Button is a telegram keyboard button
type Button struct {
	Text  string // Text is what the user sees on Button
	Value string // Value is what the send to chat Button
	Next  int32
}

// NewButton returns new Button and nil, if args is invalid, otherwise nil and errors.
func NewButton(text string, value string, next int32) (*Button, []error) {
	var err error
	errs := make([]error, 0)

	err = validateButtonText(text)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateButtonValue(value)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateButtonNextModule(next)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Button{
		Text:  text,
		Value: value,
		Next:  next,
	}, nil
}

// validateButtonText returns an error if the button text is invalid, otherwise returns nil.
func validateButtonText(text string) error {
	if len(text) == 0 {
		return ErrEmptyButtonText
	}
	return nil
}

// validateButtonValue returns an error if the button value is invalid, otherwise returns nil.
func validateButtonValue(value string) error {
	if len(value) == 0 {
		return ErrEmptyButtonValue
	}
	return nil
}

// validateButtonNextModule returns an error if the button value is invalid; otherwise returns nil.
func validateButtonNextModule(next int32) error {
	if next == 0 {
		return ErrInvalidButtonNextModule
	}
	return nil
}
