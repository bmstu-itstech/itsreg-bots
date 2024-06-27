package module

import (
	"errors"
)

type Type int32

const (
	String Type = iota + 1
	Number
)

var (
	ErrEmptyModuleTitle  = errors.New("invalid module title; should be not empty")
	ErrEmptyModuleText   = errors.New("invalid module text; should be not empty")
	ErrInvalidModuleType = errors.New("invalid module type; should be String or Number")
)

type Module struct {
	Title    string   // Title is what the Owner sees in user answers
	Text     string   // Text is what user see
	IsSilent bool     // If IsSilent, bot is not waiting for the user's response
	Type     Type     // Type defines how the user's answer is processed
	Next     *Module  // Next is the next module if module has no buttons
	Buttons  []Button // Buttons is a telegram buttons under message
}

// New returns new Module and nil, if args is invalid; otherwise nil and errors.
func New(
	title string,
	text string,
	isSilent bool,
	typ Type,
	next *Module,
	buttons []Button,
) (*Module, []error) {
	var err error
	errs := make([]error, 0)

	err = validateTitle(title)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateText(text)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateType(typ)
	if err != nil {
		errs = append(errs, err)
	}

	if buttons == nil {
		buttons = make([]Button, 0)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Module{
		Title:    title,
		Text:     text,
		IsSilent: isSilent,
		Type:     typ,
		Next:     next,
		Buttons:  buttons,
	}, nil
}

// IsFinish returns true, if there is no next module
func (m *Module) IsFinish() bool {
	return m.Next == nil
}

// validateType returns an error if the module title is invalid, otherwise returns nil.
func validateTitle(title string) error {
	if len(title) == 0 {
		return ErrEmptyModuleTitle
	}
	return nil
}

// validateType returns an error if the module text is invalid, otherwise returns nil.
func validateText(text string) error {
	if len(text) == 0 {
		return ErrEmptyModuleText
	}
	return nil
}

// validateType returns an error if the module type is invalid, otherwise returns nil.
func validateType(typ Type) error {
	switch typ {
	case String, Number:
		return nil
	}
	return ErrInvalidModuleType
}
