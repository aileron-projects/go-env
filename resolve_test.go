package env

import (
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestResolve(t *testing.T) {
	t.Setenv("TestResolve", "test")
	t.Setenv("TestResolve_Cap", "TEST")
	testCases := map[string]struct {
		in     string
		result string
		err    error
	}{
		"case01": {"", "", expressionError("")},
		"case02": {"BAD", "", expressionError("")},
		"case03": {"${!prefix}", "", expressionError("")},
		"case04": {"${/*pattern/string}", "", expressionError("")},
		"case05": {"${parameter*other}", "", expressionError("")},
		"case06": {"${!TestResolve*}", "TestResolve TestResolve_Cap", nil},
		"case07": {"${!TestResolve@}", "TestResolve TestResolve_Cap", nil},
		"case08": {"${#TestResolve}", "4", nil},
		"case09": {"${TestResolve_UnDef:-word}", "word", nil},
		"case10": {"${TestResolve_UnDef-word}", "word", nil},
		"case11": {"${TestResolve:=word}", "test", nil},
		"case12": {"${TestResolve=word}", "test", nil},
		"case13": {"${TestResolve:?word}", "test", nil},
		"case14": {"${TestResolve?word}", "test", nil},
		"case15": {"${TestResolve:+word}", "word", nil},
		"case16": {"${TestResolve+word}", "word", nil},
		"case17": {"${TestResolve:2}", "st", nil},
		"case18": {"${TestResolve:1:2}", "es", nil},
		"case19": {"${TestResolve#t}", "est", nil},
		"case20": {"${TestResolve##t}", "est", nil},
		"case21": {"${TestResolve%t}", "tes", nil},
		"case22": {"${TestResolve%%t}", "tes", nil},
		"case23": {"${TestResolve/[t]/T}", "Test", nil},
		"case24": {"${TestResolve//[^t]/X}", "tXXt", nil},
		"case25": {"${TestResolve_Cap/#[T]/t}", "tEST", nil},
		"case26": {"${TestResolve_Cap/%[T]/t}", "TESt", nil},
		"case27": {"${TestResolve^[t]}", "Test", nil},
		"case28": {"${TestResolve^^[^t]}", "tESt", nil},
		"case29": {"${TestResolve_Cap,[T]}", "tEST", nil},
		"case30": {"${TestResolve_Cap,,[^T]}", "TesT", nil},
		"case31": {"${TestResolve@U}", "TEST", nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := Resolve(tc.in)
			tester.AssertEqual(t, tc.result, v)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
}

func TestResolveGroup1(t *testing.T) {
	t.Setenv("TestResolveGroup1", "test")
	testCases := map[string]struct {
		o      string
		result string
		err    error
	}{
		"case01": {"!TestResolve*", "TestResolveGroup1", nil},
		"case02": {"!TestResolve@", "TestResolveGroup1", nil},
		"case03": {"#TestResolveGroup1", "4", nil},
		"case04": {"", "", expressionError("")},
		"case05": {"BAD", "", expressionError("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup1(tc.o)
			tester.AssertEqual(t, tc.result, v)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
}

func TestResolveGroup2(t *testing.T) {
	t.Setenv("TestResolveGroup2", "12345678")
	testCases := map[string]struct {
		p, o   string
		result string
		err    error
	}{
		"case01": {"TestResolveGroup2_UnDef", ":-word", "word", nil},
		"case02": {"TestResolveGroup2_UnDef", "-word", "word", nil},
		"case03": {"TestResolveGroup2", ":=word", "12345678", nil},
		"case04": {"TestResolveGroup2", "=word", "12345678", nil},
		"case05": {"TestResolveGroup2", ":?word", "12345678", nil},
		"case06": {"TestResolveGroup2", "?word", "12345678", nil},
		"case07": {"TestResolveGroup2", ":+word", "word", nil},
		"case08": {"TestResolveGroup2", "+word", "word", nil},
		"case09": {"TestResolveGroup2", ":5", "678", nil},
		"case10": {"TestResolveGroup2", ":3:2", "45", nil},
		"case11": {"TestResolveGroup2", "", "", expressionError("")},
		"case12": {"TestResolveGroup2", "!BAD", "", expressionError("")},
		"case13": {"TestResolveGroup2", ":", "", expressionError("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup2(tc.p, tc.o)
			tester.AssertEqual(t, tc.result, v)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
}

func TestResolveGroup3(t *testing.T) {
	t.Setenv("TestResolveGroup3", "test")
	testCases := map[string]struct {
		p, o   string
		result string
		err    error
	}{
		"case01": {"TestResolveGroup3", "#t", "est", nil},
		"case02": {"TestResolveGroup3", "##t", "est", nil},
		"case03": {"TestResolveGroup3", "%t", "tes", nil},
		"case04": {"TestResolveGroup3", "%%t", "tes", nil},
		"case05": {"TestResolveGroup3", "", "", expressionError("")},
		"case06": {"TestResolveGroup3", "!BAD", "", expressionError("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup3(tc.p, tc.o)
			tester.AssertEqual(t, tc.result, v)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
}

func TestResolveGroup4(t *testing.T) {
	t.Setenv("TestResolveGroup4_01", "test")
	t.Setenv("TestResolveGroup4_02", "TEST")
	testCases := map[string]struct {
		p, o   string
		result string
		err    error
	}{
		"case01": {"TestResolveGroup4_01", "/[t]/T", "Test", nil},
		"case02": {"TestResolveGroup4_01", "//[^t]/X", "tXXt", nil},
		"case03": {"TestResolveGroup4_02", "/#[T]/t", "tEST", nil},
		"case04": {"TestResolveGroup4_02", "/%[T]/t", "TESt", nil},
		"case05": {"TestResolveGroup4", "", "", expressionError("")},
		"case06": {"TestResolveGroup4", "/!BAD", "", expressionError("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup4(tc.p, tc.o)
			tester.AssertEqual(t, tc.result, v)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
}

func TestResolveGroup5(t *testing.T) {
	t.Setenv("TestResolveGroup5_01", "test")
	t.Setenv("TestResolveGroup5_02", "TEST")
	testCases := map[string]struct {
		p, o   string
		result string
		err    error
	}{
		"case01": {"TestResolveGroup5_01", "^[t]", "Test", nil},
		"case02": {"TestResolveGroup5_01", "^^[^t]", "tESt", nil},
		"case03": {"TestResolveGroup5_02", ",[T]", "tEST", nil},
		"case04": {"TestResolveGroup5_02", ",,[^T]", "TesT", nil},
		"case05": {"TestResolveGroup5", "", "", expressionError("")},
		"case06": {"TestResolveGroup5", "!BAD", "", expressionError("")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup5(tc.p, tc.o)
			tester.AssertEqual(t, tc.result, v)
			tester.AssertEqualErr(t, tc.err, err)
		})
	}
}
