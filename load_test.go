package env_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/aileron-projects/go-env"
	"github.com/aileron-projects/go-tester"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		got, err := env.Load("testdata/empty.txt")
		want := map[string]string{}
		for k, v := range want {
			tester.AssertEqual(t, v, os.Getenv(k))
		}
		tester.AssertDeepEqual(t, want, got)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("non empty", func(t *testing.T) {
		got, err := env.Load("testdata/comment.txt")
		want := map[string]string{"FOO": "foo"}
		for k, v := range want {
			tester.AssertEqual(t, v, os.Getenv(k))
		}
		tester.AssertDeepEqual(t, want, got)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("multiple files", func(t *testing.T) {
		got, err := env.Load("testdata/single-line-foo.txt", "testdata/single-line-bar.txt")
		want := map[string]string{"FOO": "foo", "BAR": "bar"}
		for k, v := range want {
			tester.AssertEqual(t, v, os.Getenv(k))
		}
		tester.AssertDeepEqual(t, want, got)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("duplication", func(t *testing.T) {
		got, err := env.Load("testdata/single-line-foo.txt", "testdata/single-line-bar.txt", "testdata/single-line-override.txt")
		want := map[string]string{"FOO": "foooo", "BAR": "baaar", "BAZ": "baz"}
		for k, v := range want {
			tester.AssertEqual(t, v, os.Getenv(k))
		}
		tester.AssertDeepEqual(t, want, got)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("env subst error", func(t *testing.T) {
		got, err := env.Load("testdata/env-subst-error.txt")
		tester.AssertDeepEqual(t, nil, got)
		tester.AssertEqualErr(t, &env.Error{Type: "parse"}, err)
	})
	t.Run("read file error", func(t *testing.T) {
		got, err := env.Load("testdata/not-exist.txt")
		tester.AssertDeepEqual(t, nil, got)
		tester.AssertEqualErr(t, &env.Error{Type: "load"}, err)
	})
}

func TestLoadReader(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		f, err := os.Open("testdata/empty.txt")
		tester.AssertEqualErr(t, nil, err)
		got, err := env.LoadReaders(f)
		want := map[string]string{}
		for k, v := range want {
			tester.AssertEqual(t, v, os.Getenv(k))
		}
		tester.AssertDeepEqual(t, want, got)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("non empty", func(t *testing.T) {
		f, err := os.Open("testdata/comment.txt")
		tester.AssertEqualErr(t, nil, err)
		got, err := env.LoadReaders(f)
		want := map[string]string{"FOO": "foo"}
		for k, v := range want {
			tester.AssertEqual(t, v, os.Getenv(k))
		}
		tester.AssertDeepEqual(t, want, got)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("multiple files", func(t *testing.T) {
		f1, err := os.Open("testdata/single-line-foo.txt")
		tester.AssertEqualErr(t, nil, err)
		f2, err := os.Open("testdata/single-line-bar.txt")
		tester.AssertEqualErr(t, nil, err)
		got, err := env.LoadReaders(f1, f2)
		tester.AssertEqualErr(t, nil, err)
		want := map[string]string{"FOO": "foo", "BAR": "bar"}
		for k, v := range want {
			tester.AssertEqual(t, v, os.Getenv(k))
		}
		tester.AssertDeepEqual(t, want, got)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("duplication", func(t *testing.T) {
		f1, err := os.Open("testdata/single-line-foo.txt")
		tester.AssertEqualErr(t, nil, err)
		f2, err := os.Open("testdata/single-line-bar.txt")
		tester.AssertEqualErr(t, nil, err)
		f3, err := os.Open("testdata/single-line-override.txt")
		tester.AssertEqualErr(t, nil, err)
		got, err := env.LoadReaders(f1, f2, f3)
		tester.AssertEqualErr(t, nil, err)
		want := map[string]string{"FOO": "foooo", "BAR": "baaar", "BAZ": "baz"}
		for k, v := range want {
			tester.AssertEqual(t, v, os.Getenv(k))
		}
		tester.AssertDeepEqual(t, want, got)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("env subst error", func(t *testing.T) {
		f, err := os.Open("testdata/env-subst-error.txt")
		tester.AssertEqualErr(t, nil, err)
		got, err := env.LoadReaders(f)
		tester.AssertDeepEqual(t, nil, got)
		tester.AssertEqualErr(t, &env.Error{Type: "parse"}, err)
	})
}

func TestParse(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		txt := ``
		m, err := env.Parse([]byte(txt))
		want := map[string]string{}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("comment", func(t *testing.T) {
		txt := `
		# comment line
		FOO=foo # inline comment
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"FOO": "foo"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("export", func(t *testing.T) {
		txt := `
		export FOO=foo
		export B_A_R=bar
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"FOO": "foo", "B_A_R": "bar"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("char escape", func(t *testing.T) {
		txt := `FOO=\f\o\o`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"FOO": "foo"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("quotations", func(t *testing.T) {
		txt := `
		NONE=none
		SINGLE='single'
		DOUBLE="double"
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"NONE": "none", "SINGLE": "single", "DOUBLE": "double"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("quotations in quotations", func(t *testing.T) {
		txt := `
		SINGLE='single and "double"'
		DOUBLE="'single' and double"
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"SINGLE": "single and \"double\"", "DOUBLE": "'single' and double"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("quotations escape", func(t *testing.T) {
		txt := `
		SINGLE='single \'escape\''
		DOUBLE="double \"escape\""
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"SINGLE": "single 'escape'", "DOUBLE": "double \"escape\""}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("quotations sequence", func(t *testing.T) {
		txt := `
		SEQ1='Single'"Double"
		SEQ2='Single' and "Double"
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"SEQ1": "SingleDouble", "SEQ2": "Single and Double"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("multiline", func(t *testing.T) {
		txt := `
		MULTI1='
		line1
		line2
		'
		MULTI2="
		line1
		line2
		"
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"MULTI1": "line1line2", "MULTI2": "line1line2"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("multiline with line break", func(t *testing.T) {
		txt := `
		MULTI1='
		line1\n
		line2
		'
		MULTI2="
		line1\n
		line2
		"
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"MULTI1": "line1\nline2", "MULTI2": "line1\nline2"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("end with escape", func(t *testing.T) {
		txt := `FOO=foo\`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"FOO": "foo\\"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("duplication", func(t *testing.T) {
		txt := `
		FOO=bar
		FOO=foo
		`
		m, err := env.Parse([]byte(txt))
		want := map[string]string{"FOO": "foo"}
		tester.AssertDeepEqual(t, want, m)
		tester.AssertEqualErr(t, nil, err)
	})
	t.Run("env not found", func(t *testing.T) {
		txt := `=foo`
		m, err := env.Parse([]byte(txt))
		tester.AssertDeepEqual(t, nil, m)
		tester.AssertEqualErr(t, &env.Error{Type: "parse"}, err)
	})
	t.Run("invalid char", func(t *testing.T) {
		txt := `***=foo`
		m, err := env.Parse([]byte(txt))
		tester.AssertDeepEqual(t, nil, m)
		tester.AssertEqualErr(t, &env.Error{Type: "parse"}, err)
	})
	t.Run("invalid line format", func(t *testing.T) {
		txt := `foo`
		m, err := env.Parse([]byte(txt))
		tester.AssertDeepEqual(t, nil, m)
		tester.AssertEqualErr(t, &env.Error{Type: "parse"}, err)
	})
	t.Run("quotation not closed", func(t *testing.T) {
		txt := `
		MULTI='
		line1
		line2
		`
		m, err := env.Parse([]byte(txt))
		tester.AssertDeepEqual(t, nil, m)
		tester.AssertEqualErr(t, &env.Error{Type: "parse"}, err)
	})
	t.Run("env subst error", func(t *testing.T) {
		txt := `FOO=${!FOO}`
		m, err := env.Parse([]byte(txt))
		tester.AssertDeepEqual(t, nil, m)
		tester.AssertEqualErr(t, &env.Error{Type: "parse"}, err)
	})
}

func TestParseReader(t *testing.T) {
	t.Parallel()
	t.Run("read error", func(t *testing.T) {
		txt := `
			FOO=foo
			BAR=bar
		`
		r := tester.MaxErrorReader(bytes.NewReader([]byte(txt)), 5)
		m, err := env.ParseReader(r)
		tester.AssertDeepEqual(t, nil, m)
		tester.AssertEqualErr(t, tester.ErrMaxRead, err)
	})
}
