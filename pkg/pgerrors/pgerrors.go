package pgerrors

import (
	"errors"
	"github.com/lib/pq"
)

const (
	UniqueViolationErr = pq.ErrorCode("23505")
)

func IsUniqueViolationError(err error) bool {
	var perr *pq.Error
	if errors.As(err, &perr) {
		return perr.Code == UniqueViolationErr
	}
	return false
}
