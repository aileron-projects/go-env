package env

import "regexp"

var (
	// envExp is the regular expression that matches to environmental variables.
	// See the [Resolve] for considered patterns.
	envExp = `\$\{(` +
		`[0-9a-zA-Z_]+[@]?[UuLl]?` + `|` +
		`[#!][0-9a-zA-Z_]+[*@]?` + `|` +
		`[0-9a-zA-Z_]+[:\-=?+#%/,^][^\$]*` +
		`)\}`
	envRe = regexp.MustCompilePOSIX(envExp)
)

// Subst substitutes environmental variables in the given bytes.
// See the [Resolve] for available variable syntax.
// Subst does not support nested variables.
// Use [Subst2] to allow 2 levels nested variable.
// Note that escaping variable like '\${FOO}' is not supported.
func Subst(b []byte) (sb []byte, err error) {
	defer func() {
		recover()
	}()
	sb = envRe.ReplaceAllFunc(b, func(b []byte) []byte {
		var s string
		s, err = Resolve(string(b))
		if err != nil {
			sb = nil
			panic(err)
		}
		return []byte(s)
	})
	return sb, nil
}

// Subst2 substitute environmental variable in the given bytes.
// See the [Resolve] for available variable syntax.
// Subst does support nested variables up to 2 levels.
// ${FOO_${BAR}} is allowed but ${FOO_${BAR_${BAZ}}}} is not allowed.
// Note that escaping variable like '\${FOO}' is not supported.
func Subst2(b []byte) (sb []byte, err error) {
	defer func() {
		recover()
	}()
	sb = envRe.ReplaceAllFunc(b, func(b []byte) []byte {
		var s string
		s, err = Resolve(string(b))
		if err != nil {
			sb = nil
			panic(err)
		}
		return []byte(s)
	})
	sb = envRe.ReplaceAllFunc(sb, func(b []byte) []byte {
		var s string
		s, err = Resolve(string(b))
		if err != nil {
			sb = nil
			panic(err)
		}
		return []byte(s)
	})
	return sb, nil
}
