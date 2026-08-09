package env

import (
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestSubst(t *testing.T) {
	t.Setenv("TestSubst", "test")
	t.Run("resolve all", func(t *testing.T) {
		txt := `${TestSubst}`
		b, err := Subst([]byte(txt))
		tester.AssertEqual(t, "test", string(b))
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("env error", func(t *testing.T) {
		txt := `${!TestSubst}`
		b, err := Subst([]byte(txt))
		tester.AssertEqual(t, "", string(b))
		tester.AssertEqualErr(t, &Error{Type: "expression"}, err)
	})
}

func TestSubst2(t *testing.T) {
	t.Setenv("TestSubst2", "test")
	t.Setenv("TestSubst2_PTR", "TestSubst2")
	t.Run("resolve all", func(t *testing.T) {
		txt := `${${TestSubst2_PTR}}`
		b, err := Subst2([]byte(txt))
		tester.AssertEqual(t, "test", string(b))
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("inner env error", func(t *testing.T) {
		txt := `${${!TestSubst2_PTR}}`
		b, err := Subst2([]byte(txt))
		tester.AssertEqual(t, "", string(b))
		tester.AssertEqualErr(t, &Error{Type: "expression"}, err)
	})
	t.Run("outer env error", func(t *testing.T) {
		txt := `${!${TestSubst2_PTR}}`
		b, err := Subst2([]byte(txt))
		tester.AssertEqual(t, "", string(b))
		tester.AssertEqualErr(t, &Error{Type: "expression"}, err)
	})
}
