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
	ErrEmptyModuleTitle    = errors.New("invalid module title; should be not empty")
	ErrEmptyModuleText     = errors.New("invalid module text; should be not empty")
	ErrInvalidModuleType   = errors.New("invalid module type; should be String or Number")
	ErrInvalidLastModule   = errors.New("invalid last module; if next is nil, buttons should be empty")
	ErrInvalidModuleSilent = errors.New("invalid usage silent module; should be have only default branch")
)

// Module is a message
type Module struct {
	Title    string   // Title is what the Owner sees in user answers
	Text     string   // Text is what user see
	IsSilent bool     // If IsSilent, bot is not waiting for the user's response
	Type     Type     // Type defines how the user's answer is processed
	Next     *Module  // Next is the default branch, should be nil only if the module is the last
	Buttons  []Button // Buttons are answer options
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

	if next == nil && len(buttons) > 0 {
		errs = append(errs, ErrInvalidLastModule)
	}

	if isSilent && len(buttons) > 0 {
		errs = append(errs, ErrInvalidModuleSilent)
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

func (m *Module) HasButtons() bool {
	return len(m.Buttons) > 0
}

// IsLast returns true, if there is no next module
func (m *Module) IsLast() bool {
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
