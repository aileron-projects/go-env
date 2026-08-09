package env

import (
	"cmp"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ValueType is the variable's value type.
type ValueType interface {
	~bool |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~complex64 | ~complex128 |
		~string
}

// Getenv returns value of the environmental variable.
func Getenv[T ValueType](key string) (T, error) {
	v := os.Getenv(key)
	var t T
	vv, err := parseAny(v, reflect.TypeOf(t).Kind())
	if err != nil {
		return t, &Error{Inner: err, Type: "getenv", Msg: "key=" + key}
	}
	return vv.(T), nil
}

// GetenvSlice returns values of the environmental variable.
// delim is the delimiter that separates values.
// For example, delim will be "," for values like "KEY=foo,bar,baz".
func GetenvSlice[T any](key, delim string) ([]T, error) {
	delim = cmp.Or(delim, ",")
	value := os.Getenv(key)
	arr := make([]T, 0, strings.Count(value, delim))
	if value == "" {
		return arr, nil
	}
	var t T
	for v := range strings.SplitSeq(value, delim) {
		vv, err := parseAny(v, reflect.TypeOf(t).Kind())
		if err != nil {
			return nil, &Error{Inner: err, Type: "getenv", Msg: "key=" + key}
		}
		arr = append(arr, vv.(T))
	}
	return arr, nil
}

// GetenvMap returns map data of the environmental variable.
// delim is the delimiter that separates key-value pairs.
// sep is the separator that separates key and value.
// For example, delimiter is "," and seperater is "=" for KEY="foo=alice,bar=bob".
func GetenvMap[T any](key, delim, sep string) (map[string]T, error) {
	delim = cmp.Or(delim, ",")
	sep = cmp.Or(sep, "=")
	value := os.Getenv(key)
	m := make(map[string]T, strings.Count(value, delim))
	if value == "" {
		return m, nil
	}
	var t T
	for kv := range strings.SplitSeq(value, delim) {
		before, after, _ := strings.Cut(kv, sep)
		vv, err := parseAny(after, reflect.TypeOf(t).Kind())
		if err != nil {
			return nil, &Error{Inner: err, Type: "getenv", Msg: "key=" + key}
		}
		m[before] = vv.(T)
	}
	return m, nil
}

func parseAny(v string, kind reflect.Kind) (any, error) {
	switch kind {
	case reflect.String:
		return v, nil
	case reflect.Bool:
		return strconv.ParseBool(v)
	case reflect.Int:
		return strconv.Atoi(v)
	case reflect.Int8:
		vv, err := strconv.ParseInt(v, 10, 8)
		return int8(vv), err
	case reflect.Int16:
		vv, err := strconv.ParseInt(v, 10, 16)
		return int16(vv), err
	case reflect.Int32:
		vv, err := strconv.ParseInt(v, 10, 32)
		return int32(vv), err
	case reflect.Int64:
		vv, err := strconv.ParseInt(v, 10, 64)
		return vv, err
	case reflect.Uint:
		vv, err := strconv.ParseUint(v, 10, 32)
		return uint(vv), err
	case reflect.Uint8:
		vv, err := strconv.ParseUint(v, 10, 8)
		return uint8(vv), err
	case reflect.Uint16:
		vv, err := strconv.ParseUint(v, 10, 16)
		return uint16(vv), err
	case reflect.Uint32:
		vv, err := strconv.ParseUint(v, 10, 32)
		return uint32(vv), err
	case reflect.Uint64:
		vv, err := strconv.ParseUint(v, 10, 64)
		return vv, err
	case reflect.Float32:
		vv, err := strconv.ParseFloat(v, 32)
		return float32(vv), err
	case reflect.Float64:
		vv, err := strconv.ParseFloat(v, 64)
		return vv, err
	case reflect.Complex64:
		vv, err := strconv.ParseComplex(v, 32)
		return complex64(vv), err
	case reflect.Complex128:
		vv, err := strconv.ParseComplex(v, 64)
		return vv, err
	}
	panic(fmt.Errorf("unexpected type `%s`", kind.String()))
}
