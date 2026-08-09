package env_test

import (
	"testing"

	"github.com/aileron-projects/go-env"
	"github.com/aileron-projects/go-tester"
)

func TestGetenv(t *testing.T) {
	const envName = "TestGetenv"
	t.Run("bool", func(t *testing.T) {
		t.Setenv(envName, "true")
		v, err := env.Getenv[bool](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, true, v)
	})
	t.Run("int", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[int](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("int8", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[int8](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("int16", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[int16](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("int32", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[int32](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("int64", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[int64](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("uint", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[uint](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("uint8", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[uint8](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("uint16", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[uint16](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("uint32", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[uint32](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("uint64", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[uint64](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("float32", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[float32](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("float64", func(t *testing.T) {
		t.Setenv(envName, "123")
		v, err := env.Getenv[float64](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123, v)
	})
	t.Run("complex64", func(t *testing.T) {
		t.Setenv(envName, "123+456i")
		v, err := env.Getenv[complex64](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123+456i, v)
	})
	t.Run("complex128", func(t *testing.T) {
		t.Setenv(envName, "123+456i")
		v, err := env.Getenv[complex128](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, 123+456i, v)
	})
	t.Run("string", func(t *testing.T) {
		t.Setenv(envName, "string")
		v, err := env.Getenv[string](envName)
		tester.AssertEqual(t, nil, err)
		tester.AssertEqual(t, "string", v)
	})
	t.Run("error", func(t *testing.T) {
		t.Setenv(envName, "string")
		v, err := env.Getenv[int](envName)
		tester.AssertEqualErr(t, &env.Error{Type: "getenv"}, err)
		tester.AssertEqual(t, 0, v)
	})
}

func TestGetenvSlice(t *testing.T) {
	const envName = "TestGetenvSlice"
	t.Run("empty", func(t *testing.T) {
		t.Setenv(envName, "")
		v, err := env.GetenvSlice[string](envName, "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, []string{}, v)
	})
	t.Run("single", func(t *testing.T) {
		t.Setenv(envName, "foo")
		v, err := env.GetenvSlice[string](envName, "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, []string{"foo"}, v)
	})
	t.Run("multiple", func(t *testing.T) {
		t.Setenv(envName, "foo,bar")
		v, err := env.GetenvSlice[string](envName, "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, []string{"foo", "bar"}, v)
	})
	t.Run("custom delim", func(t *testing.T) {
		t.Setenv(envName, "foo|bar")
		v, err := env.GetenvSlice[string](envName, "|")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, []string{"foo", "bar"}, v)
	})
	t.Run("type conversion", func(t *testing.T) {
		t.Setenv(envName, "true,false")
		v, err := env.GetenvSlice[bool](envName, "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, []bool{true, false}, v)
	})
	t.Run("error", func(t *testing.T) {
		t.Setenv(envName, "string")
		v, err := env.GetenvSlice[int](envName, "")
		tester.AssertEqualErr(t, &env.Error{Type: "getenv"}, err)
		tester.AssertDeepEqual(t, nil, v)
	})
}

func TestGetenvMap(t *testing.T) {
	const envName = "TestGetenvMap"
	t.Run("empty", func(t *testing.T) {
		t.Setenv(envName, "")
		v, err := env.GetenvMap[string](envName, "", "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, map[string]string{}, v)
	})
	t.Run("single", func(t *testing.T) {
		t.Setenv(envName, "foo=bar")
		v, err := env.GetenvMap[string](envName, "", "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, map[string]string{"foo": "bar"}, v)
	})
	t.Run("multiple", func(t *testing.T) {
		t.Setenv(envName, "foo=bar,alice=bob")
		v, err := env.GetenvMap[string](envName, "", "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, map[string]string{"foo": "bar", "alice": "bob"}, v)
	})
	t.Run("key only", func(t *testing.T) {
		t.Setenv(envName, "foo=bar,alice")
		v, err := env.GetenvMap[string](envName, "", "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, map[string]string{"foo": "bar", "alice": ""}, v)
	})
	t.Run("custom delim", func(t *testing.T) {
		t.Setenv(envName, "foo=bar|alice=bob")
		v, err := env.GetenvMap[string](envName, "|", "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, map[string]string{"foo": "bar", "alice": "bob"}, v)
	})
	t.Run("custom sep", func(t *testing.T) {
		t.Setenv(envName, "foo-bar,alice-bob")
		v, err := env.GetenvMap[string](envName, "", "-")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, map[string]string{"foo": "bar", "alice": "bob"}, v)
	})
	t.Run("type conversion", func(t *testing.T) {
		t.Setenv(envName, "foo=true,bar=false")
		v, err := env.GetenvMap[bool](envName, "", "")
		tester.AssertEqual(t, nil, err)
		tester.AssertDeepEqual(t, map[string]bool{"foo": true, "bar": false}, v)
	})
	t.Run("error", func(t *testing.T) {
		t.Setenv(envName, "foo=string")
		v, err := env.GetenvMap[int](envName, "", "")
		tester.AssertEqualErr(t, &env.Error{Type: "getenv"}, err)
		tester.AssertDeepEqual(t, nil, v)
	})
}
