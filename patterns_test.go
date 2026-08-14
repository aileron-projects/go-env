package env

import (
	"testing"

	"github.com/aileron-projects/go-tester"
)

// Patterns
// 	- 01: ${parameter}
// 	- 02: ${parameter:-word}
// 	- 03: ${parameter-word}
// 	- 04: ${parameter:=word}
// 	- 05: ${parameter=word}
// 	- 06: ${parameter:?word}
// 	- 07: ${parameter?word}
// 	- 08: ${parameter:+word}
// 	- 09: ${parameter+word}
// 	- 10: ${parameter:offset}
// 	- 11: ${parameter:offset:length}
// 	- 12: ${!prefix*}
// 	- 13: ${!prefix@}
// 	- 14: ${#parameter}
// 	- 15: ${parameter#word}
// 	- 16: ${parameter##word}
// 	- 17: ${parameter%word}
// 	- 18: ${parameter%%word}
// 	- 19: ${parameter/pattern/string}
// 	- 20: ${parameter//pattern/string}
// 	- 21: ${parameter/#pattern/string}
// 	- 22: ${parameter/%pattern/string}
// 	- 23: ${parameter^pattern}
// 	- 24: ${parameter^^pattern}
// 	- 25: ${parameter,pattern}
// 	- 26: ${parameter,,pattern}
// 	- 27: ${parameter@operator}

func TestEnv01(t *testing.T) {
	pattern := "${parameter}"
	testCases := map[string]struct {
		preset string
		p      string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv01", "abc", nil},
		"p set but null": {"", "TestEnv01", "", nil},
		"p not set":      {"", "TestEnv01_NotSet", "", nil},
		"invalid name":   {"", "TestEnv01--", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv01", tc.preset)
			got, err := env01(tc.p, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv01", "foo")
		got, err := env01("TestEnv01", map[string]string{"TestEnv01": "bar"})
		tester.AssertEqual(t, "bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv02(t *testing.T) {
	pattern := "${parameter:-word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv02", "word", "abc", nil},
		"p set but null": {"", "TestEnv02", "word", "word", nil},
		"p not set":      {"", "TestEnv02_NotSet", "word", "word", nil},
		"invalid name":   {"", "TestEnv--", "word", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv02", tc.preset)
			got, err := env02(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv02", "foo")
		got, err := env02("TestEnv02", "word", map[string]string{"TestEnv02": "bar"})
		tester.AssertEqual(t, "bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv03(t *testing.T) {
	pattern := "${parameter-word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv03", "word", "abc", nil},
		"p set but null": {"", "TestEnv03", "word", "", nil},
		"p not set":      {"", "TestEnv03_NotSet", "word", "word", nil},
		"invalid name":   {"", "TestEnv--", "word", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv03", tc.preset)
			got, err := env03(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv03", "foo")
		got, err := env03("TestEnv03", "word", map[string]string{"TestEnv03": "bar"})
		tester.AssertEqual(t, "bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv04(t *testing.T) {
	pattern := "${parameter:=word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv04", "word", "abc", nil},
		"p set but null": {"", "TestEnv04", "word", "word", nil},
		"p not set":      {"", "TestEnv04_NotSet", "word", "word", nil},
		"invalid name":   {"", "TestEnv--", "word", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv04", tc.preset)
			got, err := env04(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv04", "foo")
		got, err := env04("TestEnv04", "word", map[string]string{"TestEnv04": "bar"})
		tester.AssertEqual(t, "bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv05(t *testing.T) {
	pattern := "${parameter=word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv05", "word", "abc", nil},
		"p set but null": {"", "TestEnv05", "word", "", nil},
		"p not set":      {"", "TestEnv05_NotSet", "word", "word", nil},
		"invalid name":   {"", "TestEnv--", "word", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv05", tc.preset)
			got, err := env05(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv05", "foo")
		got, err := env05("TestEnv05", "word", map[string]string{"TestEnv05": "bar"})
		tester.AssertEqual(t, "bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv06(t *testing.T) {
	pattern := "${parameter:?word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv06", "word", "abc", nil},
		"p set but null": {"", "TestEnv06", "word", "", substitutionError(pattern, "", "")},
		"p not set":      {"", "TestEnv06_NotSet", "word", "", substitutionError(pattern, "", "")},
		"invalid name":   {"", "TestEnv--", "word", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv06", tc.preset)
			got, err := env06(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv06", "foo")
		got, err := env06("TestEnv06", "word", map[string]string{"TestEnv06": "bar"})
		tester.AssertEqual(t, "bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv07(t *testing.T) {
	pattern := "${parameter?word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv07", "word", "abc", nil},
		"p set but null": {"", "TestEnv07", "word", "", nil},
		"p not set":      {"", "TestEnv07_NotSet", "word", "", substitutionError(pattern, "", "")},
		"invalid name":   {"", "TestEnv--", "word", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv07", tc.preset)
			got, err := env07(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv03", "foo")
		got, err := env07("TestEnv07", "word", map[string]string{"TestEnv07": "bar"})
		tester.AssertEqual(t, "bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv08(t *testing.T) {
	pattern := "${parameter:+word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv08", "word", "word", nil},
		"p set but null": {"", "TestEnv08", "word", "", nil},
		"p not set":      {"", "TestEnv08_NotSet", "word", "", nil},
		"invalid name":   {"", "TestEnv--", "word", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv08", tc.preset)
			got, err := env08(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		got, err := env08("TestEnv08", "word", map[string]string{"TestEnv08": "bar"})
		tester.AssertEqual(t, "word", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv09(t *testing.T) {
	pattern := "${parameter+word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"p set not null": {"abc", "TestEnv09", "word", "word", nil},
		"p set but null": {"", "TestEnv09", "word", "word", nil},
		"p not set":      {"", "TestEnv09_NotSet", "word", "", nil},
		"invalid name":   {"", "TestEnv--", "word", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv09", tc.preset)
			got, err := env09(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		got, err := env09("TestEnv09", "word", map[string]string{"TestEnv09": "bar"})
		tester.AssertEqual(t, "word", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv10(t *testing.T) {
	pattern := "${parameter:offset}"
	testCases := map[string]struct {
		preset string
		p, o   string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv10", "-1", "", nil},
		"case01": {"", "TestEnv10", "0", "", nil},
		"case02": {"", "TestEnv10", "1", "", nil},
		"case03": {"a", "TestEnv10", "-1", "a", nil},
		"case04": {"a", "TestEnv10", "0", "a", nil},
		"case05": {"a", "TestEnv10", "1", "", nil},
		"case06": {"abcde", "TestEnv10", "-1", "abcde", nil},
		"case07": {"abcde", "TestEnv10", "0", "abcde", nil},
		"case08": {"abcde", "TestEnv10", "3", "de", nil},
		"case09": {"あいうえお", "TestEnv10", "3", "えお", nil},
		"case10": {"abcde", "TestEnv--", "1", "", syntaxError(nil, pattern, "", "")},
		"case11": {"abcde", "TestEnv10", "x", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv10", tc.preset)
			got, err := env10(tc.p, tc.o, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv10", "foo")
		got, err := env10("TestEnv10", "1", map[string]string{"TestEnv10": "bar"})
		tester.AssertEqual(t, "ar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv11(t *testing.T) {
	pattern := "${parameter:offset:length}"
	testCases := map[string]struct {
		preset  string
		p, o, l string
		want    string
		err     error
	}{
		"case00": {"", "TestEnv11", "-1", "0", "", nil},
		"case01": {"", "TestEnv11", "-1", "1", "", nil},
		"case02": {"", "TestEnv11", "0", "0", "", nil},
		"case03": {"", "TestEnv11", "0", "1", "", nil},
		"case04": {"a", "TestEnv11", "-1", "-1", "a", nil},
		"case05": {"a", "TestEnv11", "-1", "0", "a", nil},
		"case06": {"a", "TestEnv11", "-1", "1", "a", nil},
		"case07": {"a", "TestEnv11", "-1", "2", "a", nil},
		"case08": {"a", "TestEnv11", "0", "-1", "a", nil},
		"case09": {"a", "TestEnv11", "0", "0", "", nil},
		"case10": {"a", "TestEnv11", "0", "1", "a", nil},
		"case11": {"a", "TestEnv11", "0", "2", "a", nil},
		"case12": {"a", "TestEnv11", "1", "-1", "", nil},
		"case13": {"a", "TestEnv11", "1", "0", "", nil},
		"case14": {"a", "TestEnv11", "1", "1", "", nil},
		"case15": {"a", "TestEnv11", "1", "2", "", nil},
		"case16": {"abcde", "TestEnv11", "-1", "-1", "abcde", nil},
		"case17": {"abcde", "TestEnv11", "-1", "0", "abcde", nil},
		"case18": {"abcde", "TestEnv11", "-1", "1", "abcde", nil},
		"case19": {"abcde", "TestEnv11", "-1", "2", "abcde", nil},
		"case20": {"abcde", "TestEnv11", "0", "-1", "abcde", nil},
		"case21": {"abcde", "TestEnv11", "0", "1", "a", nil},
		"case22": {"abcde", "TestEnv11", "0", "2", "ab", nil},
		"case23": {"abcde", "TestEnv11", "0", "99", "abcde", nil},
		"case24": {"abcde", "TestEnv11", "1", "-1", "bcde", nil},
		"case25": {"abcde", "TestEnv11", "1", "0", "", nil},
		"case26": {"abcde", "TestEnv11", "1", "2", "bc", nil},
		"case27": {"abcde", "TestEnv11", "1", "99", "bcde", nil},
		"case28": {"あいうえお", "TestEnv11", "1", "2", "いう", nil},
		"case29": {"abcde", "TestEnv--", "1", "2", "", syntaxError(nil, pattern, "", "")},
		"case30": {"abcde", "TestEnv11", "x", "2", "", syntaxError(nil, pattern, "", "")},
		"case31": {"abcde", "TestEnv11", "1", "x", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv11", tc.preset)
			got, err := env11(tc.p, tc.o, tc.l, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv11", "foo")
		got, err := env11("TestEnv11", "1", "2", map[string]string{"TestEnv11": "bar"})
		tester.AssertEqual(t, "ar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv12(t *testing.T) {
	pattern := "${!prefix*}"
	testCases := map[string]struct {
		preset string
		p      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv12", "TestEnv12_var1 TestEnv12_var2", nil},
		"case01": {"a", "TestEnv12", "TestEnv12_var1 TestEnv12_var2", nil},
		"case02": {"abc", "TestEnv12", "TestEnv12_var1 TestEnv12_var2", nil},
		"case03": {"abc", "TestEnvNotFound", "", nil},
		"case04": {"abc", "TestEnv--", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv12_var1", tc.preset)
			t.Setenv("TestEnv12_var2", tc.preset)
			got, err := env12(tc.p, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv12_var1", "foo")
		got, err := env12("TestEnv12", map[string]string{"TestEnv12_var2": "bar"})
		tester.AssertEqual(t, "TestEnv12_var1 TestEnv12_var2", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv13(t *testing.T) {
	pattern := "${!prefix@}"
	testCases := map[string]struct {
		preset string
		p      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv13", "TestEnv13_var1 TestEnv13_var2", nil},
		"case01": {"a", "TestEnv13", "TestEnv13_var1 TestEnv13_var2", nil},
		"case02": {"abc", "TestEnv13", "TestEnv13_var1 TestEnv13_var2", nil},
		"case03": {"abc", "TestEnvNotFound", "", nil},
		"case04": {"abc", "TestEnv--", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv13_var1", tc.preset)
			t.Setenv("TestEnv13_var2", tc.preset)
			got, err := env13(tc.p, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv13_var1", "foo")
		got, err := env13("TestEnv13", map[string]string{"TestEnv13_var2": "bar"})
		tester.AssertEqual(t, "TestEnv13_var1 TestEnv13_var2", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv14(t *testing.T) {
	pattern := "${#parameter}"
	testCases := map[string]struct {
		preset string
		p      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv14", "0", nil},
		"case01": {"a", "TestEnv14", "1", nil},
		"case02": {"abc", "TestEnv14", "3", nil},
		"case03": {"あいうえお", "TestEnv14", "5", nil},
		"case14": {"abc", "TestEnv--", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv14", tc.preset)
			got, err := env14(tc.p, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv14", "foo")
		got, err := env14("TestEnv14", map[string]string{"TestEnv14": "barbaz"})
		tester.AssertEqual(t, "6", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv15(t *testing.T) {
	pattern := "${parameter#word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"case00": {"aabcc", "TestEnv15", "a", "abcc", nil},
		"case01": {"aabcc", "TestEnv15", "a*", "bcc", nil},
		"case02": {"aabcc", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv15", tc.preset)
			got, err := env15(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv15", "foo")
		got, err := env15("TestEnv15", "b*", map[string]string{"TestEnv15": "bar"})
		tester.AssertEqual(t, "ar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv16(t *testing.T) {
	pattern := "${parameter##word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv16", "a", "", nil},
		"case01": {"aabcc", "TestEnv16", "", "aabcc", nil},
		"case02": {"aabcc", "TestEnv16", "a", "abcc", nil},
		"case03": {"aabcc", "TestEnv16", "c", "aabcc", nil},
		"case04": {"aabcc", "TestEnv16", "^a", "abcc", nil},
		"case05": {"aabcc", "TestEnv16", "c$", "aabcc", nil},
		"case06": {"aabcc", "TestEnv16", "[ac]", "abcc", nil},
		"case07": {"aabcc", "TestEnv16", "[^ac]", "aabcc", nil},
		"case08": {"aabcc", "TestEnv16", "[a-c]", "abcc", nil},
		"case09": {"aabcc", "TestEnv16", ".", "abcc", nil},
		"case10": {"aabcc", "TestEnv16", ".?", "abcc", nil},
		"case11": {"aabcc", "TestEnv16", ".*", "", nil},
		"case12": {"aabcc", "TestEnv16", "", "aabcc", nil},
		"case13": {"aabcc", "TestEnv16", "[", "", syntaxError(nil, pattern, "", "")},
		"case14": {"aabcc", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
		"case15": {"aabcc", "TestEnv16", "a*", "bcc", nil},
		"case16": {"aabcc", "TestEnv16", "[a-b]*", "cc", nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv16", tc.preset)
			got, err := env16(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv16", "foo")
		got, err := env16("TestEnv16", "b", map[string]string{"TestEnv16": "bar"})
		tester.AssertEqual(t, "ar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv17(t *testing.T) {
	pattern := "${parameter%word}"
	testCases := map[string]struct {
		preset string
		p, w   string
		want   string
		err    error
	}{
		"case00": {"aabcc", "TestEnv17", "c", "aabc", nil},
		"case01": {"aabcc", "TestEnv17", "c*", "aab", nil},
		"case02": {"aabcc", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv17", tc.preset)
			got, err := env17(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv17", "foo")
		got, err := env17("TestEnv17", "r", map[string]string{"TestEnv17": "bar"})
		tester.AssertEqual(t, "ba", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv18(t *testing.T) {
	pattern := "${parameter%%word}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv18", "c", "", nil},
		"case01": {"aabcc", "TestEnv18", "", "aabcc", nil},
		"case02": {"aabcc", "TestEnv18", "c", "aabc", nil},
		"case03": {"aabcc", "TestEnv18", "a", "aabcc", nil},
		"case04": {"aabcc", "TestEnv18", "^a", "aabcc", nil},
		"case05": {"aabcc", "TestEnv18", "c$", "aabc", nil},
		"case06": {"aabcc", "TestEnv18", "[ac]", "aabc", nil},
		"case07": {"aabcc", "TestEnv18", "[^ac]", "aabcc", nil},
		"case08": {"aabcc", "TestEnv18", "[a-c]", "aabc", nil},
		"case09": {"aabcc", "TestEnv18", ".", "aabc", nil},
		"case10": {"aabcc", "TestEnv18", ".?", "aabc", nil},
		"case11": {"aabcc", "TestEnv18", ".*", "", nil},
		"case12": {"aabcc", "TestEnv18", "", "aabcc", nil},
		"case13": {"aabcc", "TestEnv18", "[", "", syntaxError(nil, pattern, "", "")},
		"case14": {"aabcc", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
		"case15": {"aabcc", "TestEnv18", "c*", "aab", nil},
		"case16": {"aabcc", "TestEnv18", "[b-c]*", "aa", nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv18", tc.preset)
			got, err := env18(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv18", "foo")
		got, err := env18("TestEnv18", "r", map[string]string{"TestEnv18": "bar"})
		tester.AssertEqual(t, "ba", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv19(t *testing.T) {
	pattern := "${parameter/pattern/string}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		s      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv19", "ab", "xyz", "", nil},
		"case01": {"abcabc", "TestEnv19", "", "xy", "xyabcabc", nil},
		"case02": {"abcabc", "TestEnv19", "bc", "xy", "axyabc", nil},
		"case03": {"abcabc", "TestEnv19", "ab", "xy", "xycabc", nil},
		"case04": {"abcabc", "TestEnv19", "^ab", "xy", "xycabc", nil},
		"case05": {"abcabc", "TestEnv19", "bc$", "xy", "abcaxy", nil},
		"case06": {"abcabc", "TestEnv19", "[ac]", "xy", "xybcabc", nil},
		"case07": {"abcabc", "TestEnv19", "[^ac]", "xy", "axycabc", nil},
		"case08": {"abcabc", "TestEnv19", "[a-c]", "xy", "xybcabc", nil},
		"case09": {"abcabc", "TestEnv19", ".", "xy", "xybcabc", nil},
		"case10": {"abcabc", "TestEnv19", ".?", "xy", "xybcabc", nil},
		"case11": {"abcabc", "TestEnv19", ".*", "xy", "xy", nil},
		"case12": {"abcabc", "TestEnv19", "", "", "abcabc", nil},
		"case13": {"abcabc", "TestEnv19", "[", "", "", syntaxError(nil, pattern, "", "")},
		"case14": {"abcabc", "TestEnv--", "", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv19", tc.preset)
			got, err := env19(tc.p, tc.w, tc.s, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv19", "foo")
		got, err := env19("TestEnv19", "r", "z", map[string]string{"TestEnv19": "bar"})
		tester.AssertEqual(t, "baz", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv20(t *testing.T) {
	pattern := "${parameter//pattern/string}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		s      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv20", "ab", "xyz", "", nil},
		"case01": {"abcabc", "TestEnv20", "", "x", "xaxbxcxaxbxcx", nil},
		"case02": {"abcabc", "TestEnv20", "bc", "xy", "axyaxy", nil},
		"case03": {"abcabc", "TestEnv20", "ab", "xy", "xycxyc", nil},
		"case04": {"abcabc", "TestEnv20", "^ab", "xy", "xycabc", nil},
		"case05": {"abcabc", "TestEnv20", "bc$", "xy", "abcaxy", nil},
		"case06": {"abcabc", "TestEnv20", "[ac]", "x", "xbxxbx", nil},
		"case07": {"abcabc", "TestEnv20", "[^ac]", "x", "axcaxc", nil},
		"case08": {"abcabc", "TestEnv20", "[a-c]", "x", "xxxxxx", nil},
		"case09": {"abcabc", "TestEnv20", ".", "x", "xxxxxx", nil},
		"case10": {"abcabc", "TestEnv20", ".?", "x", "xxxxxx", nil},
		"case11": {"abcabc", "TestEnv20", ".*", "xy", "xy", nil},
		"case12": {"abcabc", "TestEnv20", "", "", "abcabc", nil},
		"case13": {"abcabc", "TestEnv20", "[", "", "", syntaxError(nil, pattern, "", "")},
		"case14": {"abcabc", "TestEnv--", "", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv20", tc.preset)
			got, err := env20(tc.p, tc.w, tc.s, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv20", "foo")
		got, err := env20("TestEnv20", "r", "z", map[string]string{"TestEnv20": "bar"})
		tester.AssertEqual(t, "baz", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv21(t *testing.T) {
	pattern := "${parameter/#pattern/string}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		s      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv21", "ab", "xyz", "", nil},
		"case01": {"abcabc", "TestEnv21", "", "xy", "xyabcabc", nil},
		"case02": {"abcabc", "TestEnv21", "bc", "xy", "abcabc", nil},
		"case03": {"abcabc", "TestEnv21", "ab", "xy", "xycabc", nil},
		"case04": {"abcabc", "TestEnv21", "^ab", "xy", "xycabc", nil},
		"case05": {"abcabc", "TestEnv21", "bc$", "xy", "abcabc", nil},
		"case06": {"abcabc", "TestEnv21", "[ac]", "xy", "xybcabc", nil},
		"case07": {"abcabc", "TestEnv21", "[^ac]", "xy", "abcabc", nil},
		"case08": {"abcabc", "TestEnv21", "[a-c]", "xy", "xybcabc", nil},
		"case09": {"abcabc", "TestEnv21", ".", "xy", "xybcabc", nil},
		"case10": {"abcabc", "TestEnv21", ".?", "xy", "xybcabc", nil},
		"case11": {"abcabc", "TestEnv21", ".*", "xy", "xy", nil},
		"case12": {"abcabc", "TestEnv21", "", "", "abcabc", nil},
		"case13": {"abcabc", "TestEnv21", "[", "", "", syntaxError(nil, pattern, "", "")},
		"case14": {"abcabc", "TestEnv--", "", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv21", tc.preset)
			got, err := env21(tc.p, tc.w, tc.s, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv21", "foo")
		got, err := env21("TestEnv21", "b", "c", map[string]string{"TestEnv21": "bar"})
		tester.AssertEqual(t, "car", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv22(t *testing.T) {
	pattern := "${parameter/%pattern/string}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		s      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv22", "ab", "xyz", "", nil},
		"case01": {"abcabc", "TestEnv22", "", "xy", "abcabcxy", nil},
		"case02": {"abcabc", "TestEnv22", "bc", "xy", "abcaxy", nil},
		"case03": {"abcabc", "TestEnv22", "ab", "xy", "abcabc", nil},
		"case04": {"abcabc", "TestEnv22", "^ab", "xy", "abcabc", nil},
		"case05": {"abcabc", "TestEnv22", "bc$", "xy", "abcaxy", nil},
		"case06": {"abcabc", "TestEnv22", "[ac]", "xy", "abcabxy", nil},
		"case07": {"abcabc", "TestEnv22", "[^ac]", "xy", "abcabc", nil},
		"case08": {"abcabc", "TestEnv22", "[a-c]", "xy", "abcabxy", nil},
		"case09": {"abcabc", "TestEnv22", ".", "xy", "abcabxy", nil},
		"case10": {"abcabc", "TestEnv22", ".?", "xy", "abcabxy", nil},
		"case11": {"abcabc", "TestEnv22", ".*", "xy", "xy", nil},
		"case12": {"abcabc", "TestEnv22", "", "", "abcabc", nil},
		"case13": {"abcabc", "TestEnv22", "?", "", "", syntaxError(nil, pattern, "", "")},
		"case14": {"abcabc", "TestEnv--", "", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv22", tc.preset)
			got, err := env22(tc.p, tc.w, tc.s, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv22", "foo")
		got, err := env22("TestEnv22", "r", "z", map[string]string{"TestEnv22": "bar"})
		tester.AssertEqual(t, "baz", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv23(t *testing.T) {
	pattern := "${parameter^pattern}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv23", "a", "", nil},
		"case01": {"abcabc", "TestEnv23", "a", "Abcabc", nil},
		"case02": {"abcabc", "TestEnv23", "A", "abcabc", nil},
		"case03": {"abcabc", "TestEnv23", "ab", "abcabc", nil},
		"case04": {"abcabc", "TestEnv23", "ac", "abcabc", nil},
		"case05": {"abcabc", "TestEnv23", "^ab", "abcabc", nil},
		"case06": {"abcabc", "TestEnv23", "bc$", "abcabc", nil},
		"case07": {"abcabc", "TestEnv23", "[ac]", "Abcabc", nil},
		"case08": {"abcabc", "TestEnv23", "[^ac]", "abcabc", nil},
		"case09": {"abcabc", "TestEnv23", "[a-c]", "Abcabc", nil},
		"case10": {"abcabc", "TestEnv23", ".", "Abcabc", nil},
		"case11": {"abcabc", "TestEnv23", ".?", "Abcabc", nil},
		"case12": {"abcabc", "TestEnv23", ".*", "Abcabc", nil},
		"case13": {"abcabc", "TestEnv23", "", "Abcabc", nil},
		"case14": {"abcabc", "TestEnv23", "?", "", syntaxError(nil, pattern, "", "")},
		"case15": {"abcabc", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
		"case16": {"abcabc", "TestEnv23", "^a", "Abcabc", nil},
		"case17": {"abcabc", "TestEnv23", "a$", "Abcabc", nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv23", tc.preset)
			got, err := env23(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv23", "foo")
		got, err := env23("TestEnv23", "b", map[string]string{"TestEnv23": "bar"})
		tester.AssertEqual(t, "Bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv24(t *testing.T) {
	pattern := "${parameter^^pattern}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv24", "a", "", nil},
		"case01": {"abcabc", "TestEnv24", "a", "AbcAbc", nil},
		"case02": {"abcabc", "TestEnv24", "A", "abcabc", nil},
		"case03": {"abcabc", "TestEnv24", "ab", "ABcABc", nil},
		"case04": {"abcabc", "TestEnv24", "ac", "abcabc", nil},
		"case05": {"abcabc", "TestEnv24", "^ab", "ABcabc", nil},
		"case06": {"abcabc", "TestEnv24", "bc$", "abcaBC", nil},
		"case07": {"abcabc", "TestEnv24", "[ac]", "AbCAbC", nil},
		"case08": {"abcabc", "TestEnv24", "[^ac]", "aBcaBc", nil},
		"case09": {"abcabc", "TestEnv24", "[a-c]", "ABCABC", nil},
		"case10": {"abcabc", "TestEnv24", ".", "ABCABC", nil},
		"case11": {"abcabc", "TestEnv24", ".?", "ABCABC", nil},
		"case12": {"abcabc", "TestEnv24", ".*", "ABCABC", nil},
		"case13": {"abcabc", "TestEnv24", "", "ABCABC", nil},
		"case14": {"abcabc", "TestEnv24", "?", "", syntaxError(nil, pattern, "", "")},
		"case15": {"abcabc", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv24", tc.preset)
			got, err := env24(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv24", "foo")
		got, err := env24("TestEnv24", "b", map[string]string{"TestEnv24": "bar"})
		tester.AssertEqual(t, "Bar", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv25(t *testing.T) {
	pattern := "${parameter,pattern}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv25", "A", "", nil},
		"case01": {"ABCABC", "TestEnv25", "A", "aBCABC", nil},
		"case02": {"ABCABC", "TestEnv25", "a", "ABCABC", nil},
		"case03": {"ABCABC", "TestEnv25", "AB", "ABCABC", nil},
		"case04": {"ABCABC", "TestEnv25", "AC", "ABCABC", nil},
		"case05": {"ABCABC", "TestEnv25", "^AB", "ABCABC", nil},
		"case06": {"ABCABC", "TestEnv25", "BC$", "ABCABC", nil},
		"case07": {"ABCABC", "TestEnv25", "[AC]", "aBCABC", nil},
		"case08": {"ABCABC", "TestEnv25", "[^AC]", "ABCABC", nil},
		"case09": {"ABCABC", "TestEnv25", "[A-C]", "aBCABC", nil},
		"case10": {"ABCABC", "TestEnv25", ".", "aBCABC", nil},
		"case11": {"ABCABC", "TestEnv25", ".?", "aBCABC", nil},
		"case12": {"ABCABC", "TestEnv25", ".*", "aBCABC", nil},
		"case13": {"ABCABC", "TestEnv25", "", "aBCABC", nil},
		"case14": {"ABCABC", "TestEnv25", "?", "", syntaxError(nil, pattern, "", "")},
		"case15": {"ABCABC", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
		"case16": {"ABCABC", "TestEnv25", "^A", "aBCABC", nil},
		"case17": {"ABCABC", "TestEnv25", "A$", "aBCABC", nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv25", tc.preset)
			got, err := env25(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv25", "foo")
		got, err := env25("TestEnv25", "B", map[string]string{"TestEnv25": "BAR"})
		tester.AssertEqual(t, "bAR", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv26(t *testing.T) {
	pattern := "${parameter,,pattern}"
	testCases := map[string]struct {
		preset string
		p      string
		w      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv26", "A", "", nil},
		"case01": {"ABCABC", "TestEnv26", "A", "aBCaBC", nil},
		"case02": {"ABCABC", "TestEnv26", "a", "ABCABC", nil},
		"case03": {"ABCABC", "TestEnv26", "AB", "abCabC", nil},
		"case04": {"ABCABC", "TestEnv26", "AC", "ABCABC", nil},
		"case05": {"ABCABC", "TestEnv26", "^AB", "abCABC", nil},
		"case06": {"ABCABC", "TestEnv26", "BC$", "ABCAbc", nil},
		"case07": {"ABCABC", "TestEnv26", "[AC]", "aBcaBc", nil},
		"case08": {"ABCABC", "TestEnv26", "[^AC]", "AbCAbC", nil},
		"case09": {"ABCABC", "TestEnv26", "[A-C]", "abcabc", nil},
		"case10": {"ABCABC", "TestEnv26", ".", "abcabc", nil},
		"case11": {"ABCABC", "TestEnv26", ".?", "abcabc", nil},
		"case12": {"ABCABC", "TestEnv26", ".*", "abcabc", nil},
		"case13": {"ABCABC", "TestEnv26", "", "abcabc", nil},
		"case14": {"ABCABC", "TestEnv26", "?", "", syntaxError(nil, pattern, "", "")},
		"case15": {"ABCABC", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv26", tc.preset)
			got, err := env26(tc.p, tc.w, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv26", "foo")
		got, err := env26("TestEnv26", "B", map[string]string{"TestEnv26": "BAR"})
		tester.AssertEqual(t, "bAR", got)
		tester.AssertEqualErr(t, nil, err)
	})
}

func TestEnv27(t *testing.T) {
	pattern := "${parameter@operator}"
	testCases := map[string]struct {
		preset string
		p      string
		o      string
		want   string
		err    error
	}{
		"case00": {"", "TestEnv27", "U", "", nil},
		"case01": {"", "TestEnv27", "u", "", nil},
		"case02": {"", "TestEnv27", "L", "", nil},
		"case03": {"", "TestEnv27", "l", "", nil},
		"case04": {"abcABC", "TestEnv27", "U", "ABCABC", nil},
		"case05": {"abcABC", "TestEnv27", "u", "AbcABC", nil},
		"case06": {"ABCabc", "TestEnv27", "L", "abcabc", nil},
		"case07": {"ABCabc", "TestEnv27", "l", "aBCabc", nil},
		"case08": {"abcABC", "TestEnv27", "X", "", syntaxError(nil, pattern, "", "")},
		"case09": {"abcABC", "TestEnv--", "", "", syntaxError(nil, pattern, "", "")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TestEnv27", tc.preset)
			got, err := env27(tc.p, tc.o, nil)
			tester.AssertEqual(t, tc.want, got)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
	t.Run("priority", func(t *testing.T) {
		t.Setenv("TestEnv27", "foo")
		got, err := env27("TestEnv27", "U", map[string]string{"TestEnv27": "bar"})
		tester.AssertEqual(t, "BAR", got)
		tester.AssertEqualErr(t, nil, err)
	})
}
