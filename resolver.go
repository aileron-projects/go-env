package env

import (
	"strings"
)

// Resolve substitutes a single environmental variable.
// Supported patterns are listed below.
// Expressions are basically derived from shell parameter substitution.
// Note that the substitution behavior is NOT exactly the same as bash.
//   - https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
//   - https://tldp.org/LDP/abs/html/parameter-substitution.html
//   - https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_06_02
//
// Rules:
//
//	Expressions:
//	  01: ${parameter}                  --- See the substitution rule table below.
//	  02: ${parameter:-word}            --- See the substitution rule table below.
//	  03: ${parameter-word}             --- See the substitution rule table below.
//	  04: ${parameter:=word}            --- See the substitution rule table below.
//	  05: ${parameter=word}             --- See the substitution rule table below.
//	  06: ${parameter:?word}            --- See the substitution rule table below.
//	  07: ${parameter?word}             --- See the substitution rule table below.
//	  08: ${parameter:+word}            --- See the substitution rule table below.
//	  09: ${parameter+word}             --- See the substitution rule table below.
//	  10: ${parameter:offset}           --- Trim characters before offset.
//	  11: ${parameter:offset:length}    --- Trim characters before offset and after offset+length.
//	  12: ${!prefix*}                   --- Join the parameter name which has the prefix with a white space (Same with ${!prefix*}).
//	  13: ${!prefix@}                   --- Currently fallback to #12.
//	  14: ${#parameter}                 --- Length of value.
//	  15: ${parameter#word}             --- Currently fallback to #16.
//	  16: ${parameter##word}            --- Remove prefix of the value which matched to the word. Longest match if pattern specified.
//	  17: ${parameter%word}             --- Currently fallback to #18.
//	  18: ${parameter%%word}            --- Remove suffix of the value which matched to the word. Longest match if pattern specified.
//	  19: ${parameter/pattern/string}   --- Replace the first value which matched to the pattern to string.
//	  20: ${parameter//pattern/string}  --- Replace all values which matched to the pattern to string.
//	  21: ${parameter/#pattern/string}  --- Replace the prefix to string if matched to the pattern.
//	  22: ${parameter/%pattern/string}  --- Replace the suffix to string if matched to the pattern.
//	  23: ${parameter^pattern}          --- Convert initial character to upper case if matched to the pattern.
//	  24: ${parameter^^pattern}         --- Convert all characters which matched to the pattern to upper case.
//	  25: ${parameter,pattern}          --- Convert initial character to lower case if matched to the pattern.
//	  26: ${parameter,,pattern}         --- Convert all characters which matched to the pattern to lower case.
//	  27: ${parameter@operator}         --- Process value with the operator.
//
//	Substitution rules:
//	  |  #  |     expression     |    parameter Set     |  parameter Set  | parameter Unset |
//	  |     |                    |    and Not Null      |    But Null     |                 |
//	  | --- | ------------------ | -------------------- | --------------- | --------------- |
//	  | 01  | ${parameter}       | substitute parameter | substitute null | substitute null |
//	  | 02  | ${parameter:-word} | substitute parameter | substitute word | substitute word |
//	  | 03  | ${parameter-word}  | substitute parameter | substitute null | substitute word |
//	  | 04  | ${parameter:=word} | substitute parameter | substitute word | assign word     |
//	  | 05  | ${parameter=word}  | substitute parameter | substitute null | assign word     |
//	  | 06  | ${parameter:?word} | substitute parameter | error           | error           |
//	  | 07  | ${parameter?word}  | substitute parameter | substitute null | error           |
//	  | 08  | ${parameter:+word} | substitute word      | substitute null | substitute null |
//	  | 09  | ${parameter+word}  | substitute word      | substitute word | substitute null |
//
//	parameter:
//	  [0-9a-zA-Z_]+
//
//	word:
//	  [^\$]*
//
//	pattern:
//	  c       : matches to the character ('$' is not allowed).
//	  [a-z]   : matches specified character range.
//	  .*      : matches any length of characters.
//	  .?      : matches zero or single characters.
//
//	operator:
//	  U       : convert all characters to upper case using [strings.ToUpper]
//	  u       : convert the first character to upper case using [strings.ToUpper]
//	  L       : convert all characters to lower case using [strings.ToLower]
//	  l       : convert the first character to lower case using [strings.ToLower]
func Resolve(exp string) (string, error) {
	return resolve(exp, nil)
}

func resolve(exp string, priority map[string]string) (string, error) {
	if len(exp) < 3 || string(exp[:2]) != "${" || exp[len(exp)-1] != '}' {
		return "", expressionError(exp)
	}

	trimmed := strings.TrimSpace(exp[2 : len(exp)-1])
	parameter, others := splitVar(trimmed)
	if parameter == "" && others == "" {
		return "", expressionError(exp)
	}

	if len(others) == 0 {
		// Pattern:
		//  ${parameter}
		return env01(string(parameter), priority)
	}

	if len(parameter) == 0 {
		// Pattern:
		//  ${!prefix*} ${!prefix@}
		//  ${#parameter}
		return resolveGroup1(string(others), priority)
	}

	switch others[0] {
	case '-', '=', '?', '+':
		//  ${parameter-word} ${parameter=word}
		//  ${parameter?word} ${parameter+word}
		return resolveGroup2(string(parameter), string(others), priority)
	case ':':
		// Pattern:
		//  ${parameter:-word} ${parameter:=word}
		//  ${parameter:?word} ${parameter:+word}
		//  ${parameter:offset} ${parameter:offset:length}
		return resolveGroup2(string(parameter), string(others), priority)
	case '#', '%':
		// Pattern:
		//  ${parameter#word} ${parameter##word}
		//  ${parameter%word} ${parameter%%word}
		return resolveGroup3(string(parameter), string(others), priority)
	case '/':
		// Pattern:
		//  ${parameter/pattern/string} ${parameter//pattern/string}
		//  ${parameter/#pattern/string} ${parameter/%pattern/string}
		return resolveGroup4(string(parameter), string(others), priority)
	case '^', ',':
		// Pattern:
		//  ${parameter^pattern} ${parameter^^pattern}
		//  ${parameter,pattern} ${parameter,,pattern}
		return resolveGroup5(string(parameter), string(others), priority)
	case '@':
		// Pattern:
		//  ${parameter@operator}
		return env27(string(parameter), string(others[1:]), priority)
	default:
		return "", expressionError(exp)
	}
}

func splitVar(s string) (parameter, others string) {
	for i, c := range []byte(s) {
		if !validChar(c) {
			return s[:i], s[i:]
		}
	}
	return s, ""
}

// resolveGroup1 resolves the following pattern group.
//   - 12: ${!prefix*}   --> o=!prefix*
//   - 13: ${!prefix@}   --> o=!prefix@
//   - 14: ${#parameter} --> o=#parameter
func resolveGroup1(o string, priority map[string]string) (string, error) {
	if len(o) < 2 {
		return "", expressionError("${" + o + "}")
	}
	switch o[0] {
	case '#':
		return env14(o[1:], priority)
	case '!':
		switch o[len(o)-1] {
		case '*':
			return env12(o[1:len(o)-1], priority)
		case '@':
			return env13(o[1:len(o)-1], priority)
		}
	}
	return "", expressionError("${" + o + "}")
}

// resolveGroup2 resolves the following pattern group.
//   - 02: ${parameter:-word}         --> o=:-word
//   - 03: ${parameter-word}          --> o=-word
//   - 04: ${parameter:=word}         --> o=:=word
//   - 05: ${parameter=word}          --> o==word
//   - 06: ${parameter:?word}         --> o=:?word
//   - 07: ${parameter?word}          --> o=?word
//   - 08: ${parameter:+word}         --> o=:+word
//   - 09: ${parameter+word}          --> o=+word
//   - 10: ${parameter:offset}        --> o=:offset
//   - 11: ${parameter:offset:length} --> o=:offset:length
func resolveGroup2(p, o string, priority map[string]string) (string, error) {
	if len(o) < 1 {
		return "", expressionError("${" + p + o + "}")
	}
	switch o[0] {
	case '-':
		return env03(p, o[1:], priority)
	case '=':
		return env05(p, o[1:], priority)
	case '?':
		return env07(p, o[1:], priority)
	case '+':
		return env09(p, o[1:], priority)
	case ':':
		if len(o) < 2 {
			return "", expressionError("${" + p + o + "}")
		}
		switch o[1] {
		case '-':
			return env02(p, o[2:], priority)
		case '=':
			return env04(p, o[2:], priority)
		case '?':
			return env06(p, o[2:], priority)
		case '+':
			return env08(p, o[2:], priority)
		default:
			if i := strings.Index(o[1:], ":"); i > 0 {
				return env11(p, o[1:i+1], o[i+2:], priority)
			} else {
				return env10(p, o[1:], priority)
			}
		}
	}
	return "", expressionError("${" + p + o + "}")
}

// resolveGroup3 resolves the following pattern group.
//   - 15: ${parameter#word}   --> o=#word
//   - 16: ${parameter##word}  --> o=##word
//   - 17: ${parameter%word}   --> o=%word
//   - 18: ${parameter%%word}  --> o=%%word
func resolveGroup3(p, o string, priority map[string]string) (string, error) {
	if len(o) < 1 {
		return "", expressionError("${" + p + o + "}")
	}
	switch o[0] {
	case '#':
		if len(o) > 1 && o[1] == '#' {
			return env16(p, o[2:], priority)
		} else {
			return env15(p, o[1:], priority)
		}
	case '%':
		if len(o) > 1 && o[1] == '%' {
			return env18(p, o[2:], priority)
		} else {
			return env17(p, o[1:], priority)
		}
	}
	return "", expressionError("${" + p + o + "}")
}

// resolveGroup4 resolves the following pattern group.
//   - 19: ${parameter/pattern/string}  --> o=/pattern/string
//   - 20: ${parameter//pattern/string} --> o=//pattern/string
//   - 21: ${parameter/#pattern/string} --> o=/#pattern/string
//   - 22: ${parameter/%pattern/string} --> o=/%pattern/string
func resolveGroup4(p, o string, priority map[string]string) (string, error) {
	if len(o) < 2 {
		return "", expressionError("${" + p + o + "}")
	}
	switch o[1] {
	case '/':
		if i := strings.Index(o[2:], "/"); i > 0 {
			return env20(p, o[2:i+2], o[i+3:], priority)
		}
	case '#':
		if i := strings.Index(o[2:], "/"); i > 0 {
			return env21(p, o[2:i+2], o[i+3:], priority)
		}
	case '%':
		if i := strings.Index(o[2:], "/"); i > 0 {
			return env22(p, o[2:i+2], o[i+3:], priority)
		}
	default:
		if i := strings.Index(o[1:], "/"); i > 0 {
			return env19(p, o[1:i+1], o[i+2:], priority)
		}
	}
	return "", expressionError("${" + p + o + "}")
}

// resolveGroup5 resolves the following pattern group.
//   - 23: ${parameter^pattern}  --> o=^pattern
//   - 24: ${parameter^^pattern} --> o=^^pattern
//   - 25: ${parameter,pattern}  --> o=,pattern
//   - 26: ${parameter,,pattern} --> o=,,pattern
func resolveGroup5(p, o string, priority map[string]string) (string, error) {
	if len(o) < 1 {
		return "", expressionError("${" + p + o + "}")
	}
	switch o[0] {
	case '^':
		if len(o) > 1 && o[1] == '^' {
			return env24(p, o[2:], priority)
		} else {
			return env23(p, o[1:], priority)
		}
	case ',':
		if len(o) > 1 && o[1] == ',' {
			return env26(p, o[2:], priority)
		} else {
			return env25(p, o[1:], priority)
		}
	}
	return "", expressionError("${" + p + o + "}")
}
