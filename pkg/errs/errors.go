package errs

import (
	"errors"
	"fmt"
	"runtime"
)

var ErrNotFound = errors.New("not found")

type cartError struct {
	msg      string
	file     string
	line     int
	fn       string
	err      error
	HTTPCode int
}

func (e *cartError) Error() string {
	return fmt.Sprintf("%s:%v", e.msg, e.err)
}

func New(msg string) error {
	e := &cartError{msg: msg}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}
	e.file = file
	e.line = line

	e.fn = runtime.FuncForPC(pc).Name()

	return e
}

func NewWithError(msg string, err error) error {
	e := &cartError{msg: msg, err: err}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return e
	}
	e.file = file
	e.line = line

	e.fn = runtime.FuncForPC(pc).Name()

	return e
}

// Unwrap is used to make it work with errors.Is, errors.As.
func (e *cartError) Unwrap() error {
	return e.err
}

func (e *cartError) Func() string {
	return e.fn
}

func (e *cartError) Filepath() string {
	return fmt.Sprintf("%s:%d", e.file, e.line)
}

func (e *cartError) LogFields() map[string]string {
	return map[string]string{
		"error": e.Error(),
		"file":  e.Filepath(),
		"func":  e.Func(),
	}
}
