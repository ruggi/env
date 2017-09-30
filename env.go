package env

import (
	"os"
	"reflect"
	"strconv"
)

const DefaultTag = "env"

// ParseInto is the same as ParseTagInto, but using the DefaultTag tag.
func ParseInto(c interface{}) {
	ParseTagInto(DefaultTag, c)
}

// ParseTagInto tries to set field values in an interface using environment variables.
// The environment variables must have the same name as the one specified in the field's
// tag (with the same tag name passed as argument to ParseTagInto).
// If the interface argument is not a pointer, ParseTagInto returns immediately.
// The following types are supported: string, int, int8, int32, int64, float32, float64, bool.
// If ParseTagInto fails converting a field, it just skips to the next field in the struct.
func ParseTagInto(tag string, c interface{}) {
	ParseTagIntoFunc(tag, c, defaultConv)
}

// ParseTagIntoFunc performs the same as ParseTagInto, but instead of using the default conversion
// functions for environment variables, uses the provided ConvFunc when a tag-environment variable match
// is found.
func ParseTagIntoFunc(tag string, c interface{}, convFn ConvFunc) {
	if reflect.TypeOf(c).Kind() != reflect.Ptr {
		return
	}
	typ := reflect.TypeOf(c).Elem()
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.Field(i)
		lookup, ok := fld.Tag.Lookup(tag)
		if !ok {
			continue
		}
		env := os.Getenv(lookup)
		if env == "" {
			continue
		}
		val := reflect.ValueOf(c).Elem().Field(i)
		dest, ok := convFn(env, val.Type().Kind())
		if !ok {
			continue
		}
		val.Set(reflect.ValueOf(dest))
	}
}

// ConvFunc is a function used for converting environment variables into fields.
type ConvFunc func(envVar string, kind reflect.Kind) (interface{}, bool)

type singleConvFunc func(envVar string) (interface{}, error)

var convs = map[reflect.Kind]singleConvFunc{
	reflect.String:  func(v string) (interface{}, error) { return v, nil },
	reflect.Int8:    int8Conv,
	reflect.Int16:   int16Conv,
	reflect.Int:     intConv,
	reflect.Int32:   int32Conv,
	reflect.Int64:   int64Conv,
	reflect.Float32: float32Conv,
	reflect.Float64: float64Conv,
	reflect.Bool:    boolConv,
}

func stringConv(v string) (interface{}, error) {
	return v, nil
}

func int8Conv(v string) (interface{}, error) {
	i, err := strconv.ParseInt(v, 10, 8)
	return int8(i), err
}

func int16Conv(v string) (interface{}, error) {
	i, err := strconv.ParseInt(v, 10, 16)
	return int16(i), err
}

func intConv(v string) (interface{}, error) {
	return strconv.Atoi(v)
}

func int32Conv(v string) (interface{}, error) {
	i, err := strconv.ParseInt(v, 10, 32)
	return int32(i), err
}

func int64Conv(v string) (interface{}, error) {
	return strconv.ParseInt(v, 10, 64)
}

func float32Conv(v string) (interface{}, error) {
	f, err := strconv.ParseFloat(v, 32)
	return float32(f), err
}

func float64Conv(v string) (interface{}, error) {
	return strconv.ParseFloat(v, 64)
}

func boolConv(v string) (interface{}, error) {
	return strconv.ParseBool(v)
}

func defaultConv(v string, k reflect.Kind) (interface{}, bool) {
	fn, ok := convs[k]
	if !ok {
		return nil, false
	}
	conv, err := fn(v)
	if err != nil {
		return nil, false
	}
	return conv, true
}
