package env

import (
	"regexp"
)

var (
	// envPattern = `(` +
	// 	`[0-9a-zA-Z_]+[@]?[UuLl]?` + `|` +
	// 	`[#!][0-9a-zA-Z_]+[*@]?` + `|` +
	// 	`[0-9a-zA-Z_]+[:\-=?+#%/,^][^\$]*` + `)`
	envRe    = regexp.MustCompile(`\$\{\s*[^{}]*?\s*\}`)
	envEscRe = regexp.MustCompile(`\\\$\\\{\s*.*?\s*\\\}`)
)

// Subst substitute environmental variable in the given bytes.
// See the [Resolve] for available variable syntax.
// Subst supports nested variables like ${FOO_${BAR}}.
// Use '\\' to escape the expression. e.g. \$\{FOO\}.
func Subst(b []byte) ([]byte, error) {
	return subst(b, nil)
}

func subst(b []byte, priority map[string]string) (sb []byte, err error) {
	defer func() {
		recover()
	}()
	sb = b
	left := 10 // Max repeat.
	found := true
	for left > 0 && found {
		left -= 1
		found = false
		sb = envRe.ReplaceAllFunc(sb, func(b []byte) []byte {
			found = true
			s, e := resolve(string(b), priority)
			if e != nil {
				sb = nil
				err = e
				panic(e)
			}
			return []byte(s)
		})
	}
	left = 10 // Max repeat.
	found = true
	for left > 0 && found {
		left -= 1
		found = false
		sb = envEscRe.ReplaceAllFunc(sb, func(b []byte) []byte {
			found = true
			b = b[2 : len(b)-1] // Reuse the buffer.
			b[0], b[1], b[len(b)-1] = '$', '{', '}'
			return b
		})
	}
	return sb, nil
}
