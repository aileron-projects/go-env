package env

import (
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestSubst(t *testing.T) {
	t.Setenv("TestSubst", "test")
	t.Setenv("TestSubst_PTR", "TestSubst")
	t.Run("single", func(t *testing.T) {
		txt := `\$\{TestSubst\}=${TestSubst}`
		b, err := Subst([]byte(txt))
		tester.AssertEqual(t, "${TestSubst}=test", string(b))
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("nested", func(t *testing.T) {
		txt := `\$\{\$\{TestSubst_PTR\}\}=${${TestSubst_PTR}}`
		b, err := Subst([]byte(txt))
		tester.AssertEqual(t, "${${TestSubst_PTR}}=test", string(b))
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("inner env error", func(t *testing.T) {
		txt := `${${!TestSubst_PTR}}`
		b, err := Subst([]byte(txt))
		tester.AssertEqual(t, "", string(b))
		tester.AssertEqualErr(t, &Error{Type: "expression"}, err)
	})
	t.Run("outer env error", func(t *testing.T) {
		txt := `${!${TestSubst_PTR}}`
		b, err := Subst([]byte(txt))
		tester.AssertEqual(t, "", string(b))
		tester.AssertEqualErr(t, &Error{Type: "expression"}, err)
	})
}
