package bot

import (
	"errors"
	"github.com/zhikh23/itsreg-tgservice/internal/domain/module"
)

var (
	ErrInvalidOwnerId     = errors.New("invalid owner id; should be > 0")
	ErrEmptyName          = errors.New("invalid name; should be not empty")
	ErrEmptyTgToken       = errors.New("invalid telegram token; should be not empty")
	ErrEmptyGsToken       = errors.New("invalid google sheets token; should be not empty")
	ErrInvalidLimit       = errors.New("invalid limit; should be > 0")
	ErrInvalidStartModule = errors.New("invalid start module; should not nil")
)

// Bot describes the settings of the telegram bot script.
type Bot struct {
	OwnerId int32
	Name    string
	TgToken string
	GsToken string
	Limit   int32 // Maximum members in command; 1, if a single event
	Start   *module.Module
}

// New returns new Bot and nil, if args is invalid, otherwise nil and errors.
func New(
	ownerId int32,
	name string,
	tgToken string,
	gsToken string,
	teamLimit int32,
	start *module.Module,
) (*Bot, []error) {
	var err error
	errs := make([]error, 0)

	err = validateOwnerId(ownerId)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateName(name)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateTgToken(tgToken)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateGsToken(gsToken)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateLimit(teamLimit)
	if err != nil {
		errs = append(errs, err)
	}

	err = validateStartModule(start)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Bot{
		OwnerId: ownerId,
		Name:    name,
		TgToken: tgToken,
		GsToken: gsToken,
		Limit:   teamLimit,
		Start:   start,
	}, nil
}

// validateOwnerId returns an error if the bot's owner id is invalid, otherwise returns nil.
func validateOwnerId(id int32) error {
	if id <= 0 {
		return ErrInvalidOwnerId
	}
	return nil
}

// validateName returns an error if the bot's test is invalid, otherwise returns nil.
func validateName(name string) error {
	if len(name) == 0 {
		return ErrEmptyName
	}
	return nil
}

// validateTgToken returns an error if the telegram bot token is invalid, otherwise returns nil.
func validateTgToken(tok string) error {
	if len(tok) == 0 {
		return ErrEmptyTgToken
	}
	return nil
}

// validateGsToken returns an error if the bot's GoogleSheets token is invalid, otherwise returns nil.
func validateGsToken(tok string) error {
	if len(tok) == 0 {
		return ErrEmptyGsToken
	}
	return nil
}

// validateLimitToken returns an error if the bot's limit is invalid, otherwise returns nil.
func validateLimit(limit int32) error {
	if limit <= 0 {
		return ErrInvalidLimit
	}
	return nil
}

// validateStartModule return an error, if the bot's start module is invalid; otherwise returns nil.
func validateStartModule(module *module.Module) error {
	if module == nil {
		return ErrInvalidStartModule
	}

	// TODO: checking for recursion?

	return nil
}
