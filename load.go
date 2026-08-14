package env

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
)

// Load loads environmental variables from files.
// It sets parsed values with [os.Setenv].
// Duplicated keys are always overwritten.
// The default ".env" is loaded if no files provided.
// See [Parse] for file formats and [Resolve] for enviromental variable expressions.
func Load(files ...string) (map[string]string, error) {
	if len(files) == 0 {
		files = append(files, ".env")
	}
	readers := make([]io.Reader, 0, len(files))
	for _, f := range files {
		r, err := os.Open(f)
		if err != nil {
			return nil, &Error{Inner: err, Type: "load", Msg: f}
		}
		readers = append(readers, r)
		defer r.Close()
	}
	return LoadReaders(readers...)
}

// LoadReaders loads environmental variables from files.
// It sets parsed values with [os.Setenv].
// Duplicated keys are always overwritten.
// Unline [Load], LoadReaders do nothing even when no readers were provided.
func LoadReaders(readers ...io.Reader) (map[string]string, error) {
	m := map[string]string{}
	for _, r := range readers {
		envs, err := parse(r)
		if err != nil {
			return nil, err
		}
		for k, v := range envs {
			if err := os.Setenv(k, v); err != nil {
				return nil, &Error{Inner: err, Type: "load", Msg: "set env `" + k + "`"}
			}
			m[k] = v
		}
	}
	return m, nil
}

// Parse parses environmental variable from the given bytes.
// Typically Parse parses variables from files such as .env file.
// Parse resolves embedded environmental variables in the b.
// See [Resolve] for expressions.
//
// References:
//   - https://github.com/joho/godotenv
//   - https://github.com/motdotla/dotenv
//
// Input specifications:
//
//	Single line:
//		# Single quotes and double quotes are removed if entire value is enclosed.
//		# "export" can be placed before env name.
//		FOO=BAR          >> BAR
//		FOO="BAR"        >> BAR
//		FOO='BAR'        >> BAR
//		FOO='B"R'        >> B"R
//		FOO="B'R"        >> B'R
//		export FOO=BAR   >> BAR
//
//	Multiple lines:
//		# The following definition of FOO results in "BARBAZ".
//		# Line breaks of LF and CRLF are removed.
//		# BOTH single quotes and double quotes can be used to enclose multiple lines.
//		FOO="
//		BAR
//		BAZ
//		"
//
//	Comments:
//		# Sharp '#' can be used for commenting.
//		# It must not be in the scope of single quotes and double quotes.
//		# It must have at least 1 white space before '#' if the comment is inlined.
//		# comment            >> Comment is appropriately parsed.
//		FOO=BAR # comment    >> Comment is appropriately parsed.
//		FOO=BAR# comment     >> '#' is not parsed as comment. It considered as a part of value.
//
//	Escapes:
//		# '\\' can be used for escaping characters by following the 3 rules.
//		# 1. '\\' always escapes special character of ', ", \\, #
//		# 2. '\\' is ignored when it is not in the scope of single quotes or double quotes.
//		# 3. '\\'n or "\n" in the scope of single or doubles quotes results in line breaks of LF.
//		FOO=B\"R      >> B"R
//		FOO=B\'R      >> B'A
//		FOO="B\"R"    >> B"R
//		FOO=B\R       >> BR (Its not in a scope of single or double quotes.)
//		FOO="B\nR"    >> B<LF>R (\n is, if in a scope of quotes, converted into a line break.)
//
//	Environmental variables:
//		# Load resolves environmental variables.
//		FOO=BAR${BAZ}
func Parse(b []byte) (map[string]string, error) {
	return ParseReader(bytes.NewReader(b))
}

// Parse parses environmental variables from the reader.
// See [Parse] for more details.
func ParseReader(r io.Reader) (map[string]string, error) {
	envs, err := parse(r)
	if err != nil {
		return nil, err
	}
	return envs, nil
}

func parse(r io.Reader) (map[string]string, error) {
	reader := bufio.NewReader(r)
	envs := map[string]string{}
	inSingleQuote, inDoubleQuote := false, false
	multilineName, multilineValue := "", ""
	current := 0 // Line number.
	eof := false // End of file flag.
	for !eof {
		current += 1
		line, err := reader.ReadBytes('\n')
		eof = (err == io.EOF)
		if !eof && err != nil {
			return nil, err
		}
		line = bytes.Trim(line, "\t\n\f\r ")
		line, err = subst(line, envs) // Replace environmental variable if exists.
		if err != nil {
			return nil, &Error{Inner: err, Type: "parse", Msg: "line " + strconv.Itoa(current)}
		}

		if inSingleQuote || inDoubleQuote {
			var val string
			val, inSingleQuote, inDoubleQuote = scanValue(line, inSingleQuote, inDoubleQuote)
			multilineValue += val
			if !inSingleQuote && !inDoubleQuote {
				envs[multilineName] = multilineValue
				multilineName = ""  // Reset variable.
				multilineValue = "" // Reset variable.
			}
			if eof {
				err := errors.New("quotation not closed `" + multilineName + "`")
				return nil, &Error{Inner: err, Type: "parse", Msg: "line " + strconv.Itoa(current)}
			}
		} else {
			line = bytes.TrimPrefix(line, []byte("export"))
			line = bytes.TrimLeft(line, "\t ")
			name, rest, err := scanName(line)
			if err != nil {
				return nil, &Error{Inner: err, Type: "parse", Msg: "line " + strconv.Itoa(current)}
			}
			if name == "" {
				continue // Maybe comment line.
			}
			var val string
			val, inSingleQuote, inDoubleQuote = scanValue(rest, inSingleQuote, inDoubleQuote)
			if inSingleQuote || inDoubleQuote {
				multilineName = name
				multilineValue = val
			} else {
				envs[name] = val
			}
		}
	}
	return envs, nil
}

// scanKey scans a line and looks for environmental variable key name.
// If the line is comment, it returns empty key.
// If a key found, it returns the key and the rest of the line after '='.
//
// Example input patterns:
//   - FOO=bar
//   - FOO="bar"
//   - FOO="bar" #comment
//   - FOO="bar
//   - #comment
func scanName(b []byte) (name string, rest []byte, err error) {
	if len(b) == 0 || b[0] == '#' {
		return "", nil, nil
	}
	if b[0] == '=' {
		return "", nil, errors.New("key not found")
	}
	for i, c := range b {
		switch {
		case '0' <= c && c <= '9':
		case 'a' <= c && c <= 'z':
		case 'A' <= c && c <= 'Z':
		case c == '_':
		case c == '=':
			return string(b[:i]), b[i+1:], nil
		default:
			return "", nil, errors.New("invalid character `" + string(c) + "`")
		}
	}
	return "", nil, errors.New("invalid expression")
}

// scanValue scans value line.
// It returns the value and the flag of
// in-single-quote or in-double-quote.
// isq and idq don't become true simultaneously.
//
// Example input patterns:
//
//	foobar        <-- Non quoted value
//	"foobar"      <-- Quoted value
//	"foobar       <-- Quoted value. Double quote is not closed.
//	'foobar       <-- Quoted value. Single quote is not closed.
//	"foo"'bar     <-- Quoted value. Single quote is not closed.
//	"foo"'bar'    <-- Non quoted value. All quotes are closed.
//	foo#comment   <-- Non quoted value. The '#' is treated as a value.
//	foo #comment  <-- Commented value. Comments are ignored.
//	#comment      <-- Commented value. Comments are ignored.
func scanValue(b []byte, sq, dq bool) (val string, isq, idq bool) {
	v := make([]byte, 0, len(b))
	escaped := false
	inSingleQuote := sq
	inDoubleQuote := dq
	for _, c := range b {
		if escaped {
			if c == '\\' || c == '\'' || c == '"' || c == '#' {
				v = append(v, c)
				escaped = false
				continue
			}
			if (inSingleQuote || inDoubleQuote) && c == 'n' {
				v = append(v, '\n')
				escaped = false
				continue
			}
			v = append(v, c)
			escaped = false
			continue
		}
		switch c {
		case '\\':
			escaped = true
		case '\'':
			if inDoubleQuote {
				v = append(v, c)
			} else {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if inSingleQuote {
				v = append(v, c)
			} else {
				inDoubleQuote = !inDoubleQuote
			}
		case '#':
			if !inSingleQuote && !inDoubleQuote && (len(v) > 1 && v[len(v)-1] == ' ') {
				v = bytes.TrimRight(v, " ")
			}
			return string(v), inSingleQuote, inDoubleQuote
		default:
			v = append(v, c)
		}
	}
	if escaped {
		v = append(v, '\\')
	}
	return string(v), inSingleQuote, inDoubleQuote
}
