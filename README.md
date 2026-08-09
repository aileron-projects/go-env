# go-env

**Go library for dealing with environmental variables.**

<div align="center">

[![GoDoc](https://godoc.org/github.com/aileron-projects/go-env?status.svg)](http://godoc.org/github.com/aileron-projects/go-env)
[![Test](https://github.com/aileron-projects/go-env/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/aileron-projects/go-env/actions/workflows/test.yaml?query=branch%3Amain)
[![License](https://img.shields.io/badge/License-Apache%202.0-yellow.svg)](./LICENSE)

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-env)
[![OpenSourceInsight](https://badgen.net/badge/open%2Fsource%2F/insight/cyan)](https://deps.dev/go/github.com%2Faileron-projects%2Fgo-env)
[![OSS Insight](https://badgen.net/badge/OSS/Insight/orange)](https://ossinsight.io/analyze/aileron-projects/go-env)

</div>

## Features

- Get environmental variables with generics function
- Resolve environmental variable with bash-like expressions.
- Loading environmental variable files.

## Tested Environments

Operating System:

- `Linux`: [ubuntu-latest](https://github.com/actions/runner-images)
- `Windows`: [windows-latest](https://github.com/actions/runner-images)
- `macOS`: [macos-latest](https://github.com/actions/runner-images)

Architecture (Using QEMU on linux):

- x86: `amd64`, `386`
- arm: `arm/v5`, `arm/v6`, `arm/v7`, `arm64`
- risc: `riscv64`, `loong64`
- ppc: `ppc64`, `ppc64le`
- mips: `mips`, `mips64`, `mips64le`, `mipsle`
- ibm: `s390x`

## Release Cycle

- Releases are made as needed.
- [Semantic Versioning](https://semver.org/) `vX.Y.Z` is used.

## License

[Apache-2.0](LICENSE)

## Usage

### Resolving and substituting environment variables

Single variable expression can be resolved with `env.Resolve`.
Use `env.Subst` or `env.Subst2` to resolve variables in text contents.
`env.Subst2` can resolve nested environment like `${${FOO}}`.

```go
// Use Resolve to resolve single variable.
value, err := env.Resolve("${FOO}")
value, err := env.Resolve("${BAR:-default}")

// Use Subst or Subst2 to resolve variables in texts.
txt := `
FOO is ${FOO}.
BAR is ${BAR:-default}.
`
result, err := env.Subst([]byte(txt))
```

Supported expressions:

01. `${parameter}`                  --- See the substitution rule table below.
02. `${parameter:-word}`            --- See the substitution rule table below.
03. `${parameter-word}`             --- See the substitution rule table below.
04. `${parameter:=word}`            --- See the substitution rule table below.
05. `${parameter=word}`             --- See the substitution rule table below.
06. `${parameter:?word}`            --- See the substitution rule table below.
07. `${parameter?word}`             --- See the substitution rule table below.
08. `${parameter:+word}`            --- See the substitution rule table below.
09. `${parameter+word}`             --- See the substitution rule table below.
10. `${parameter:offset}`           --- Trim characters before offset.
11. `${parameter:offset:length}`    --- Trim characters before offset and after offset+length.
12. `${!prefix*}`                   --- Join the parameter name which has the prefix with a white space (Same with ${!prefix*}).
13. `${!prefix@}`                   --- Currently fallback to #12.
14. `${#parameter}`                 --- Length of value.
15. `${parameter#word}`             --- Currently fallback to #16.
16. `${parameter##word}`            --- Remove prefix of the value which matched to the word. Longest match if pattern specified.
17. `${parameter%word}`             --- Currently fallback to #18.
18. `${parameter%%word}`            --- Remove suffix of the value which matched to the word. Longest match if pattern specified.
19. `${parameter/pattern/string}`   --- Replace the first value which matched to the pattern to string.
20. `${parameter//pattern/string}`  --- Replace all values which matched to the pattern to string.
21. `${parameter/#pattern/string}`  --- Replace the prefix to string if matched to the pattern.
22. `${parameter/%pattern/string}`  --- Replace the suffix to string if matched to the pattern.
23. `${parameter^pattern}`          --- Convert initial character to upper case if matched to the pattern.
24. `${parameter^^pattern}`         --- Convert all characters which matched to the pattern to upper case.
25. `${parameter,pattern}`          --- Convert initial character to lower case if matched to the pattern.
26. `${parameter,,pattern}`         --- Convert all characters which matched to the pattern to lower case.
27. `${parameter@operator}`         --- Process value with the operator.

Substitution rules:

| #   | expression         | parameter Set and Not Null | parameter Set but Null | parameter Unset |
| --- | ------------------ | -------------------------- | ---------------------- | --------------- |
| 01  | ${parameter}       | substitute parameter       | substitute null        | substitute null |
| 02  | ${parameter:-word} | substitute parameter       | substitute word        | substitute word |
| 03  | ${parameter-word}  | substitute parameter       | substitute null        | substitute word |
| 04  | ${parameter:=word} | substitute parameter       | substitute word        | assign word     |
| 05  | ${parameter=word}  | substitute parameter       | substitute null        | assign word     |
| 06  | ${parameter:?word} | substitute parameter       | error                  | error           |
| 07  | ${parameter?word}  | substitute parameter       | substitute null        | error           |
| 08  | ${parameter:+word} | substitute word            | substitute null        | substitute null |
| 09  | ${parameter+word}  | substitute word            | substitute word        | substitute null |

```txt
parameter:
  [0-9a-zA-Z_]+

word:
  [^\$]*

pattern:
  c       : matches to the character ('$' is not allowed).
  [a-z]   : matches specified character range.
  .*      : matches any length of characters.
  .?      : matches zero or single characters.

operator:
  U       : convert all characters to upper case using [strings.ToUpper]
  u       : convert the first character to upper case using [strings.ToUpper]
  L       : convert all characters to lower case using [strings.ToLower]
  l       : convert the first character to lower case using [strings.ToLower]
```

### Loading env files

- `env.Load` loads environmental variables from files and set values by `os.Setenv`
- `env.LoadReaders` works like env.Load but it takes io.Reader instead
- `env.Parse` parses environmental variables without calling `os.Setenv`
- `env.ParseReader` works like env.Parse but it takes io.Reader instead

```go
err := env.Load()                        // Loads ".env"
err := env.Load("prod.env")              // Loads custom file
err := env.Load("common.env", "dev.env") // Loads multiple files
```

Environmental variable files can be written in the following formats.

**Single line:**

```txt
# Single quotes and double quotes are removed if entire value is enclosed.
# "export" can be placed before name.
FOO=BAR          # BAR
FOO="BAR"        # BAR
FOO='BAR'        # BAR
FOO='B"R'        # B"R
FOO="B'R"        # B'R
export FOO=BAR   # BAR
```

**Multiple lines:**

```txt
# The following definition of FOO results in "BARBAZ".
# Line breaks of LF and CRLF are removed.
# BOTH single quotes and double quotes can be used to enclose multiple lines.
FOO="
BAR
BAZ
"
```

**Comments:**

```txt
# Sharp '#' can be used for commenting.
# It must not be in the scope of single quotes and double quotes.
# It must have at least 1 white space before '#' if the comment is inlined.
# comment            # Comment is appropriately parsed.
FOO=BAR # comment    # Comment is appropriately parsed.
FOO=BAR# comment     # '#' is not parsed as comment. It considered as a part of value.
```

**Escapes:**

```txt
# '\\' can be used for escaping characters by following the 3 rules.
# 1. '\\' always escapes special character of ', ", \\, #
# 2. '\\' is ignored when it is not in the scope of single quotes or double quotes.
# 3. '\\'n or "\n" in the scope of single or doubles quotes results in line breaks of LF.
FOO=B\"R     # B"R
FOO=B\'R     # B'A
FOO="B\"R"   # B"R
FOO=B\R      # BR (Its not in a scope of single or double quotes.)
FOO="B\nR"   # B<LF>R (\n is, if in a scope of quotes, converted into a line break.)
```

**Environmental variables:**

```txt
# Load resolves environmental variables.
FOO=${BAR}
```

### Auto-loading env file

Import `autoload` package to automatically load `.env`.
File path can be changed by `autoload.FilePath`.

```go
import (
    _ "github.com/aileron-projects/go-env/autoload"
)
```

### Getting environmental variable values

`env.Getenv` can be used for getting a single value.

Supported types are `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `complex64`, `complex128` and `string`.

```go
os.Setenv("FOO", "foo")
v, err := env.Getenv[string]("FOO") // "foo"

os.Setenv("BAR", "123")
v, err := env.Getenv[int]("BAR") // 123

os.Setenv("BAZ", "true")
v, err := env.Getenv[bool]("BAZ") // true
```

Slices and maps are also supported by `env.GetenvSlice` and `env.GetenvMap` each.

```go
os.Setenv("FOO", "alice,bob")
s, err := env.GetenvSlice[string]("FOO", ",")  // [alice bob]

os.Setenv("BAR", "alice|bob")
s, err := env.GetenvSlice[string]("BAR", "|") // [alice bob]

os.Setenv("BAZ", "123,456")
s, err := env.GetenvSlice[int]("BAZ", "") // [123 456]
```

```go
os.Setenv("FOO", "key1=val1,key2=val2")
m, err := env.GetenvMap[string]("FOO", "", "") // map[key1:val1 key2:val2]

os.Setenv("BAR", "key1:val1|key2:val2|key3")
m, err := env.GetenvMap[string]("BAR", "|", ":") // map[key1:val1 key2:val2 key3:]

os.Setenv("BAZ", "key1=123,key2=456")
m, err := env.GetenvMap[int]("BAZ", "", "") // map[key1:123 key2:456]
```

## Build Tags

No build tags defined for this library.

## Enviromental Variables

No environmental variables defined for this library.
