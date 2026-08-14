package env

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestError(t *testing.T) {
	t.Parallel()
	t.Run("unwrap", func(t *testing.T) {
		err := &Error{Inner: io.EOF}
		inner := err.Unwrap()
		tester.AssertEqualErr(t, io.EOF, inner)
	})
	t.Run("error message", func(t *testing.T) {
		err := &Error{Inner: io.EOF, Type: "type", Msg: "msg"}
		msg := err.Error()
		tester.AssertEqual(t, "go-env/env: type: msg [EOF]", msg)
	})
	t.Run("empty message", func(t *testing.T) {
		err := &Error{Inner: io.EOF, Type: "type", Msg: ""}
		msg := err.Error()
		tester.AssertEqual(t, "go-env/env: type: [EOF]", msg)
	})
	t.Run("errors equal", func(t *testing.T) {
		err1 := &Error{Type: "type", Msg: "aaa", Inner: nil}
		err2 := &Error{Type: "type", Msg: "bbb", Inner: io.EOF}
		tester.AssertEqualErr(t, err1, err2)
	})
	t.Run("wrapped error equal", func(t *testing.T) {
		err1 := &Error{Type: "type", Msg: "aaa", Inner: nil}
		err2 := &Error{Type: "type", Msg: "bbb", Inner: io.EOF}
		err3 := fmt.Errorf("outer error [%w]", err2)
		tester.AssertEqual(t, true, err1.Is(err3))
		tester.AssertEqualErr(t, err1, err3)
	})
	t.Run("errors not equal", func(t *testing.T) {
		err1 := &Error{Type: "foo", Msg: "", Inner: nil}
		err2 := &Error{Type: "bar", Msg: "", Inner: io.EOF}
		tester.AssertEqual(t, false, errors.Is(err1, err2))
	})
	t.Run("wrapped error not equal", func(t *testing.T) {
		err1 := &Error{Type: "type", Msg: "msg", Inner: nil}
		err2 := fmt.Errorf("outer error [%w]", io.EOF)
		tester.AssertEqual(t, false, err1.Is(err2))
	})
}
