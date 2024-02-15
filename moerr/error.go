package moerr

import (
	"errors"
	"fmt"
)

func NewInternalErrorNoCtx(msg string, args ...any) error {
	xmsg := fmt.Sprintf(msg, args...)
	return errors.New(xmsg)
}

func NewDivByZeroNoCtx() error {
	return errors.New("division by zero")
}

func NewArrayInvalidOpNoCtx(expected, actual int) error {
	xmsg := fmt.Sprintf("vector ops between different dimensions (%v, %v) is not permitted.", expected, actual)
	return errors.New(xmsg)
}
