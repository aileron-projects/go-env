package env_test

import (
	"fmt"
	"os"

	"github.com/aileron-projects/go-env"
)

func ExampleGetenv() {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "123")
	os.Setenv("BAZ", "true")
	os.Setenv("ERR", "string")

	fmt.Println(env.Getenv[string]("FOO"))
	fmt.Println(env.Getenv[int]("BAR"))
	fmt.Println(env.Getenv[bool]("BAZ"))
	fmt.Println(env.Getenv[int]("ERR"))
	// Output:
	// foo <nil>
	// 123 <nil>
	// true <nil>
	// 0 go-env/env: getenv: key=ERR [strconv.Atoi: parsing "string": invalid syntax]
}

func ExampleGetenvSlice() {
	os.Setenv("FOO", "alice,bob")
	os.Setenv("BAR", "alice|bob")
	os.Setenv("BAZ", "123,456")
	os.Setenv("ERR", "string")

	fmt.Println(env.GetenvSlice[string]("FOO", ",")) // Default delim is ","
	fmt.Println(env.GetenvSlice[string]("BAR", "|")) // Custom delimiter
	fmt.Println(env.GetenvSlice[int]("BAZ", ""))
	fmt.Println(env.GetenvSlice[int]("ERR", ""))
	// Output:
	// [alice bob] <nil>
	// [alice bob] <nil>
	// [123 456] <nil>
	// [] go-env/env: getenv: key=ERR [strconv.Atoi: parsing "string": invalid syntax]
}

func ExampleGetenvMap() {
	os.Setenv("FOO", "key1=val1,key2=val2")
	os.Setenv("BAR", "key1:val1|key2:val2|key3")
	os.Setenv("BAZ", "key1=123,key2=456")
	os.Setenv("ERR", "key=string")

	fmt.Println(env.GetenvMap[string]("FOO", "", ""))   // Default delimiter "," and seperator "="
	fmt.Println(env.GetenvMap[string]("BAR", "|", ":")) // Custom delimiter and seperator
	fmt.Println(env.GetenvMap[int]("BAZ", "", ""))
	fmt.Println(env.GetenvMap[int]("ERR", "", ""))
	// Output:
	// map[key1:val1 key2:val2] <nil>
	// map[key1:val1 key2:val2 key3:] <nil>
	// map[key1:123 key2:456] <nil>
	// map[] go-env/env: getenv: key=ERR [strconv.Atoi: parsing "string": invalid syntax]
}

func ExampleSubst() {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "FOO")

	b1, _ := env.Subst([]byte(`{FOO}=${FOO}`))
	fmt.Println(string(b1))

	b2, _ := env.Subst([]byte(`{{BAR}}=${${BAR}}`))
	fmt.Println(string(b2)) // Nested env is not supported.

	// Output:
	// {FOO}=foo
	// {{BAR}}=${FOO}
}

func ExampleSubst2() {
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "FOO")

	b1, _ := env.Subst2([]byte(`{FOO}=${FOO}`))
	fmt.Println(string(b1))

	b2, _ := env.Subst2([]byte(`{{BAR}}={FOO}=${${BAR}}`))
	fmt.Println(string(b2)) // Nested env is not supported.

	// Output:
	// {FOO}=foo
	// {{BAR}}={FOO}=foo
}

func ExampleSubst_all() {
	os.Setenv("ABC", "abcdefg")
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "BAR")
	os.Setenv("ARR_X", "xxx")
	os.Setenv("ARR_Y", "yyy")
	os.Unsetenv("BAZ")

	txt := []byte(`
01: {parameter}                 => ${FOO}
02: {parameter:-word}           => ${BAZ:-default}
03: {parameter-word}            => ${BAZ-default}
04: {parameter:=word}           => ${BAZ:=default}
05: {parameter=word}            => ${BAZ=default}
06: {parameter:?word}           => ${BAZ:?default}
07: {parameter?word}            => ${BAZ?default}
08: {parameter:+word}           => ${BAZ:+default}
09: {parameter+word}            => ${BAZ+default}
10: {parameter:offset}          => ${ABC:3}
11: {parameter:offset:length}   => ${ABC:3:3}
12: {!prefix*}                  => ${!ARR*}
13: {!prefix@}                  => ${!ARR@}
14: {#parameter}                => ${#FOO}
15: {parameter#word}            => ${FOO#[a-z]}
16: {parameter##word}           => ${FOO##[a-z]}
17: {parameter%word}            => ${FOO%[a-z]}
18: {parameter%%word}           => ${FOO%%[a-z]}
19: {parameter/pattern/string}  => ${FOO/[a-z]/x}
20: {parameter//pattern/string} => ${FOO//[a-z]/x}
21: {parameter/#pattern/string} => ${FOO/#[a-z]/x}
22: {parameter/%pattern/string} => ${FOO/%[a-z]/x}
23: {parameter^pattern}         => ${FOO^[f]}
24: {parameter^^pattern}        => ${FOO^^[o]}
25: {parameter,pattern}         => ${BAR,[B]}
26: {parameter,,pattern}        => ${BAR,,[A]}
27: {parameter@U}               => ${FOO@U}
27: {parameter@u}               => ${FOO@u}
27: {parameter@L}               => ${BAR@L}
27: {parameter@l}               => ${BAR@l}
`)

	b, _ := env.Subst(txt)
	fmt.Println(string(b))
	// Output:
	// 01: {parameter}                 => foo
	// 02: {parameter:-word}           => default
	// 03: {parameter-word}            => default
	// 04: {parameter:=word}           => default
	// 05: {parameter=word}            => default
	// 06: {parameter:?word}           => default
	// 07: {parameter?word}            => default
	// 08: {parameter:+word}           => default
	// 09: {parameter+word}            => default
	// 10: {parameter:offset}          => defg
	// 11: {parameter:offset:length}   => def
	// 12: {!prefix*}                  => ARR_X ARR_Y
	// 13: {!prefix@}                  => ARR_X ARR_Y
	// 14: {#parameter}                => 3
	// 15: {parameter#word}            => oo
	// 16: {parameter##word}           => oo
	// 17: {parameter%word}            => fo
	// 18: {parameter%%word}           => fo
	// 19: {parameter/pattern/string}  => xoo
	// 20: {parameter//pattern/string} => xxx
	// 21: {parameter/#pattern/string} => xoo
	// 22: {parameter/%pattern/string} => fox
	// 23: {parameter^pattern}         => Foo
	// 24: {parameter^^pattern}        => fOO
	// 25: {parameter,pattern}         => bAR
	// 26: {parameter,,pattern}        => BaR
	// 27: {parameter@U}               => FOO
	// 27: {parameter@u}               => Foo
	// 27: {parameter@L}               => bar
	// 27: {parameter@l}               => bAR
}

func ExampleResolve() {
	os.Setenv("ABC", "abcdefg")
	os.Setenv("FOO", "foo")
	os.Setenv("BAR", "BAR")
	os.Setenv("ARR_X", "xxx")
	os.Setenv("ARR_Y", "yyy")
	os.Unsetenv("BAZ")

	must := func(s string, err error) string {
		if err != nil {
			panic(err)
		}
		return s
	}
	fmt.Println("${FOO} ------------", must(env.Resolve("${FOO}")))
	fmt.Println("${BAZ:-default} ---", must(env.Resolve("${BAZ:-default}")))
	fmt.Println("${BAZ-default}  ---", must(env.Resolve("${BAZ-default}")))
	fmt.Println("${BAZ:=default} ---", must(env.Resolve("${BAZ:=default}")))
	fmt.Println("${BAZ=default}  ---", must(env.Resolve("${BAZ=default}")))
	fmt.Println("${BAZ:?default} ---", must(env.Resolve("${BAZ:?default}")))
	fmt.Println("${BAZ?default}  ---", must(env.Resolve("${BAZ?default}")))
	fmt.Println("${BAZ:+default} ---", must(env.Resolve("${BAZ:+default}")))
	fmt.Println("${BAZ+default}  ---", must(env.Resolve("${BAZ+default}")))
	fmt.Println("${ABC:3} ----------", must(env.Resolve("${ABC:3}")))
	fmt.Println("${ABC:3:3} --------", must(env.Resolve("${ABC:3:3}")))
	fmt.Println("${!ARR*} ----------", must(env.Resolve("${!ARR*}")))
	fmt.Println("${!ARR@} ----------", must(env.Resolve("${!ARR@}")))
	fmt.Println("${#FOO} ----------", must(env.Resolve("${#FOO}")))
	fmt.Println("${FOO#[a-z]} -----", must(env.Resolve("${FOO#[a-z]}")))
	fmt.Println("${FOO##[a-z]} ----", must(env.Resolve("${FOO##[a-z]}")))
	fmt.Println("${FOO%[a-z]} -----", must(env.Resolve("${FOO%[a-z]}")))
	fmt.Println("${FOO%%[a-z]} ----", must(env.Resolve("${FOO%%[a-z]}")))
	fmt.Println("${FOO/[a-z]/x} ---", must(env.Resolve("${FOO/[a-z]/x}")))
	fmt.Println("${FOO//[a-z]/x} --", must(env.Resolve("${FOO//[a-z]/x}")))
	fmt.Println("${FOO/#[a-z]/x} --", must(env.Resolve("${FOO/#[a-z]/x}")))
	fmt.Println("${FOO/%[a-z]/x} --", must(env.Resolve("${FOO/%[a-z]/x}")))
	fmt.Println("${FOO^[f]} -------", must(env.Resolve("${FOO^[f]}")))
	fmt.Println("${FOO^^[o]} ------", must(env.Resolve("${FOO^^[o]}")))
	fmt.Println("${BAR,[B]} -------", must(env.Resolve("${BAR,[B]}")))
	fmt.Println("${BAR,,[A]} ------", must(env.Resolve("${BAR,,[A]}")))
	fmt.Println("${FOO@U} ---------", must(env.Resolve("${FOO@U}")))
	// Output:
	// ${FOO} ------------ foo
	// ${BAZ:-default} --- default
	// ${BAZ-default}  --- default
	// ${BAZ:=default} --- default
	// ${BAZ=default}  --- default
	// ${BAZ:?default} --- default
	// ${BAZ?default}  --- default
	// ${BAZ:+default} --- default
	// ${BAZ+default}  --- default
	// ${ABC:3} ---------- defg
	// ${ABC:3:3} -------- def
	// ${!ARR*} ---------- ARR_X ARR_Y
	// ${!ARR@} ---------- ARR_X ARR_Y
	// ${#FOO} ---------- 3
	// ${FOO#[a-z]} ----- oo
	// ${FOO##[a-z]} ---- oo
	// ${FOO%[a-z]} ----- fo
	// ${FOO%%[a-z]} ---- fo
	// ${FOO/[a-z]/x} --- xoo
	// ${FOO//[a-z]/x} -- xxx
	// ${FOO/#[a-z]/x} -- xoo
	// ${FOO/%[a-z]/x} -- fox
	// ${FOO^[f]} ------- Foo
	// ${FOO^^[o]} ------ fOO
	// ${BAR,[B]} ------- bAR
	// ${BAR,,[A]} ------ BaR
	// ${FOO@U} --------- FOO
}
