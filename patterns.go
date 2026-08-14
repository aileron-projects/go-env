package env

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// lookupEnv retrieves the value of the environment variable named by the key.
// It returns false and empty string if no values found for the key.
// It checks the priority map and returns the value if found
// before calling [os.LookupEnv].
func lookupEnv(key string, priority map[string]string) (string, bool) {
	if priority != nil {
		if v, ok := priority[key]; ok {
			return v, true
		}
	}
	return os.LookupEnv(key)
}

// validParamName returns if the parameter name, or enviromental name
// is valid pattern or not.
func validParamName(s string) error {
	for i := 0; i < len(s); i++ {
		if !validChar(s[i]) {
			return fmt.Errorf("character `%s` is not allowed", string(s[i]))
		}
	}
	return nil
}

func validChar(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'z':
		return true
	case 'A' <= c && c <= 'Z':
		return true
	case c == '_':
		return true
	default:
		return false
	}
}

// env01 returns the value of:
//   - 01: ${parameter}
func env01(p string, priority map[string]string) (string, error) {
	const pattern = "${parameter}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, _ := lookupEnv(p, priority)
	return v, nil
}

// env02 returns the value of:
//   - 02: ${parameter:-word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: substitute word
//   - parameter Unset: substitute word
func env02(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter:-word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, ok := lookupEnv(p, priority)
	if !ok { // parameter unset
		return w, nil
	}
	if v == "" { // parameter set but null
		return w, nil
	}
	return v, nil // parameter set not null
}

// env03 returns the value of:
//   - 03: ${parameter-word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: substitute null
//   - parameter Unset: substitute word
func env03(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter-word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, ok := lookupEnv(p, priority)
	if !ok { // parameter unset
		return w, nil
	}
	if v == "" { // parameter set but null
		return "", nil
	}
	return v, nil // parameter set not null
}

// env04 returns the value of:
//   - 04: ${parameter:=word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: assign word
//   - parameter Unset: assign word
func env04(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter:=word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, ok := lookupEnv(p, priority)
	if !ok { // parameter unset
		return w, os.Setenv(p, w)
	}
	if v == "" { // parameter set but null
		return w, os.Setenv(p, w)
	}
	return v, nil // parameter set not null
}

// env05 returns the value of:
//   - 05: ${parameter=word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: substitute null
//   - parameter Unset: assign word
func env05(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter=word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, ok := lookupEnv(p, priority)
	if !ok { // parameter unset
		return w, os.Setenv(p, w)
	}
	if v == "" { // parameter set but null
		return "", nil
	}
	return v, nil // parameter set not null
}

// env06 returns the value of:
//   - 06: ${parameter:?word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: error
//   - parameter Unset: error
func env06(p, _ string, priority map[string]string) (string, error) {
	const pattern = "${parameter:?word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, ok := lookupEnv(p, priority)
	if !ok {
		return "", substitutionError(pattern, p, "parameter unset")
	}
	if v == "" {
		return "", substitutionError(pattern, p, "parameter set but null")
	}
	return v, nil // parameter set not empty
}

// env07 returns the value of:
//   - 07: ${parameter?word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: substitute null
//   - parameter Unset: error
func env07(p, _ string, priority map[string]string) (string, error) {
	const pattern = "${parameter?word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, ok := lookupEnv(p, priority)
	if !ok {
		return "", substitutionError(pattern, p, "parameter unset")
	}
	if v == "" { // parameter set but null
		return "", nil
	}
	return v, nil // parameter set not null
}

// env08 returns the value of:
//   - 08: ${parameter:+word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute word
//   - parameter Set But Null: substitute null
//   - parameter Unset: substitute null
func env08(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter:+word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, ok := lookupEnv(p, priority)
	if !ok { // parameter is not set
		return "", nil
	}
	if v == "" { // parameter set but null
		return "", nil
	}
	return w, nil // parameter set not null
}

// env09 returns the value of:
//   - 09: ${parameter+word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute word
//   - parameter Set But Null: substitute word
//   - parameter Unset: substitute null
func env09(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter+word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, ok := lookupEnv(p, priority)
	if !ok { // parameter unset
		return "", nil
	}
	if v == "" { // parameter set but null
		return w, nil
	}
	return w, nil // parameter set not null
}

// env10 returns the value of:
//   - 10: ${parameter:offset}
func env10(p, o string, priority map[string]string) (string, error) {
	const pattern = "${parameter:offset}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	offset, err := strconv.Atoi(o)
	if err != nil {
		return "", syntaxError(err, pattern, p, "parseing offset")
	}
	v, _ := lookupEnv(p, priority)
	r := []rune(v)
	if offset < 0 {
		return v, nil
	}
	if offset > len(r) {
		return "", nil
	}
	return string(r[offset:]), nil
}

// env11 returns the value of:
//   - 11: ${parameter:offset:length}
func env11(p, o, l string, priority map[string]string) (string, error) {
	const pattern = "${parameter:offset:length}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	offset, err := strconv.Atoi(o)
	if err != nil {
		return "", syntaxError(err, pattern, p, "parseing offset")
	}
	length, err := strconv.Atoi(l)
	if err != nil {
		return "", syntaxError(err, pattern, p, "parseing length")
	}
	v, _ := lookupEnv(p, priority)
	if offset < 0 {
		return v, nil
	}
	r := []rune(v)
	if offset >= len(r) {
		return "", nil
	}
	if length < 0 {
		return string(r[offset:]), nil
	}
	if offset+length > len(r) {
		return string(r[offset:]), nil
	}
	return string(r[offset : offset+length]), nil
}

// env12 returns the value of:
//   - 12: ${!prefix*}
func env12(p string, priority map[string]string) (string, error) {
	const pattern = "${!prefix*}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	names := []string{}
	for k, _ := range priority {
		if strings.HasPrefix(k, p) {
			names = append(names, k)
		}
	}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, p) {
			names = append(names, strings.Split(v, "=")[0])
		}
	}
	slices.Sort(names)
	names = slices.Compact(names) // Remove duplications.
	return strings.Join(names, " "), nil
}

// env13 returns the value of:
//   - 13: ${!prefix@}
func env13(p string, priority map[string]string) (string, error) {
	const pattern = "${!prefix@}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	return env12(p, priority) // Fallback
}

// env14 returns the value of:
//   - 14: ${#parameter}
func env14(p string, priority map[string]string) (string, error) {
	const pattern = "${#parameter}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, _ := lookupEnv(p, priority)
	return strconv.Itoa(len([]rune(v))), nil
}

// env15 returns the value of:
//   - 15: ${parameter#word}
func env15(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter#word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	return env16(p, w, priority) // Fallback.
}

// env16 returns the value of:
//   - 16: ${parameter##word}
func env16(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter##word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX("^" + w)
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	i := re.FindAllStringIndex(v, -1)
	if len(i) == 0 {
		return v, nil
	}
	return v[i[len(i)-1][1]:], nil
}

// env17 returns the value of:
//   - 17: ${parameter%word}
func env17(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter%word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	return env18(p, w, priority) // Fallback.
}

// env18 returns the value of:
//   - 18: ${parameter%%word}
func env18(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter%%word}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX(w + "$")
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	i := re.FindAllStringIndex(v, -1)
	if len(i) == 0 {
		return v, nil
	}
	return v[:i[len(i)-1][0]], nil // Remove longest match.
}

// env19 returns the value of:
//   - 19: ${parameter/pattern/string}
func env19(p, w, s string, priority map[string]string) (string, error) {
	const pattern = "${parameter/pattern/string}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX(w)
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	replaced := false
	v = re.ReplaceAllStringFunc(v, func(ss string) string {
		if replaced {
			return ss
		}
		replaced = true
		return s
	})
	return v, nil
}

// env20 returns the value of:
//   - 20: ${parameter//pattern/string}
func env20(p, w, s string, priority map[string]string) (string, error) {
	const pattern = "${parameter//pattern/string}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX(w)
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	return re.ReplaceAllString(v, s), nil
}

// env21 returns the value of:
//   - 21: ${parameter/#pattern/string}
func env21(p, w, s string, priority map[string]string) (string, error) {
	const pattern = "${parameter/#pattern/string}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX("^" + w)
	if err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, _ := lookupEnv(p, priority)
	i := re.FindAllStringIndex(v, -1)
	if len(i) == 0 {
		return v, nil
	}
	return s + v[i[len(i)-1][1]:], nil
}

// env22 returns the value of:
//   - 22: ${parameter/%pattern/string}
func env22(p, w, s string, priority map[string]string) (string, error) {
	const pattern = "${parameter/%pattern/string}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX(w + "$")
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	i := re.FindAllStringIndex(v, -1)
	if len(i) == 0 {
		return v, nil
	}
	return v[:i[len(i)-1][0]] + s, nil
}

// env23 returns the value of:
//   - 23: ${parameter^pattern}
func env23(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter^pattern}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX(cmp.Or(w, ".?"))
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	if len(v) == 0 {
		return "", nil
	}
	return re.ReplaceAllStringFunc(v[:1], strings.ToUpper) + v[1:], nil
}

// env24 returns the value of:
//   - 24: ${parameter^^pattern}
func env24(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter^^pattern}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX(cmp.Or(w, ".?"))
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	return re.ReplaceAllStringFunc(v, strings.ToUpper), nil
}

// env25 returns the value of:
//   - 25: ${parameter,pattern}
func env25(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter,pattern}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX(cmp.Or(w, ".?"))
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	if len(v) == 0 {
		return "", nil
	}
	return re.ReplaceAllStringFunc(v[:1], strings.ToLower) + v[1:], nil
}

// env26 returns the value of:
//   - 26: ${parameter,,pattern}
func env26(p, w string, priority map[string]string) (string, error) {
	const pattern = "${parameter,,pattern}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	re, err := regexp.CompilePOSIX(cmp.Or(w, ".?"))
	if err != nil {
		return "", syntaxError(err, pattern, p, "expression `"+w+"`")
	}
	v, _ := lookupEnv(p, priority)
	return re.ReplaceAllStringFunc(v, strings.ToLower), nil
}

// env27 returns the value of:
//   - 27: ${parameter@operator}
func env27(p, o string, priority map[string]string) (string, error) {
	const pattern = "${parameter@operator}"
	if err := validParamName(p); err != nil {
		return "", syntaxError(err, pattern, p, "")
	}
	v, _ := lookupEnv(p, priority)
	if v == "" {
		return "", nil
	}
	switch o {
	case "U":
		return strings.ToUpper(v), nil
	case "u":
		return strings.ToUpper(v[:1]) + v[1:], nil
	case "L":
		return strings.ToLower(v), nil
	case "l":
		return strings.ToLower(v[:1]) + v[1:], nil
	}
	err := fmt.Errorf("unsupported operator `%s`", o)
	return "", syntaxError(err, pattern, p, "")
}
