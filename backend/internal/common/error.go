package common

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportedLang = fmt.Errorf("unsupported language")
	ErrInvalidLevelId  = fmt.Errorf("invalid level id")
	ErrEmptyInput      = fmt.Errorf("empty input")
	ErrBadSubmit       = fmt.Errorf("bad submit")
)

func FromErrors(baseErr error, targetErrs ...error) bool {
	for _, err := range targetErrs {
		if errors.Is(baseErr, err) {
			return true
		}
	}
	return false
}

func IsValidateErr(err error) bool {
	return FromErrors(err, ErrUnsupportedLang, ErrEmptyInput, ErrInvalidLevelId)
}

func NewBadSubmitErr(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrBadSubmit)
}
