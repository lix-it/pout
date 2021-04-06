package pout

import (
	"errors"
	"fmt"
)

var (
	ErrTrailingDot = errors.New("trailing dot after message name")
)

type ErrInvalidIdentifier struct {
	Detail error
}

func (e ErrInvalidIdentifier) Error() string {
	return fmt.Sprintf("invalid package/message identifier given: %s", e.Detail.Error())
}

func (e ErrInvalidIdentifier) Unwrap() error {
	return e.Detail
}
