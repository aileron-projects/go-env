package env

import (
	"errors"
)

// Error is environmental variable related error.
type Error struct {
	Inner error  // Inner is the inner error.
	Type  string // Type is the error type.
	Msg   string // Msg is the error message.
}

func (e *Error) Unwrap() error {
	return e.Inner
}

func (e *Error) Error() string {
	s := "go-env/env: " + e.Type + ":"
	if e.Msg != "" {
		s += " " + e.Msg
	}
	if e.Inner != nil {
		s = s + " [" + e.Inner.Error() + "]"
	}
	return s
}

func (e *Error) Is(err error) bool {
	for err != nil {
		ee, ok := err.(*Error)
		if ok {
			return e.Type == ee.Type
		}
		err = errors.Unwrap(err)
	}
	return false
}

func expressionError(expression string) *Error {
	return &Error{
		Type: "expression",
		Msg:  "invalid `" + expression + "`",
	}
}

func syntaxError(inner error, pattern, parameter, explain string) *Error {
	detail := "parameter `" + parameter + "`"
	if explain != "" {
		detail = detail + ". " + explain
	}
	return &Error{
		Inner: inner,
		Type:  "syntax",
		Msg:   "pattern " + pattern + ". " + detail,
	}
}

func substitutionError(pattern, parameter, explain string) *Error {
	detail := "parameter `" + parameter + "`"
	if explain != "" {
		detail = detail + ". " + explain
	}
	return &Error{
		Type: "substitution",
		Msg:  "pattern " + pattern + ". " + detail,
	}
}
