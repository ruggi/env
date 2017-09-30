# env

[![godoc](https://godoc.org/github.com/ruggi/env?status.svg)](https://godoc.org/github.com/ruggi/env)

env tries to set a struct's field values with environment variables' values.

It comes in handy when dealing with configuration types, inits, tests, benchmarks, etc.

## Usage

The typical use case is when you have a type with default values but you want to customize some of them at run time without the need for additional code.

With this prerequisite:

```go
type Server struct {
    Address string `env:"SERVER_ADDR"`
    Port    int    `env:"SERVER_PORT"` 
}

func init() {
    srv := &Server{
        Address: "localhost",
        Port:    8080,
    }
    env.ParseInto(srv)
}
```

You can override some (or all) values of `Server` with environment variables named as the `env` tags in the struct's declaration.

```
$ SERVER_PORT=4567 go run example.go
```

When calling `env.ParseInto` (or the other parsing methods), the interface argument must be a pointer. If it's not, the method returns immediately.

Supported types:
* string
* int, int8, int16, int32, int64
* float32, float64
* bool

If the conversion of an environment variable fails, it just skips the field without changing the original value.

You can support more types with your own conversion function using the `env.ParseTagIntoFunc` method.

### Methods

```go
func ParseInto(c interface{})
```
Uses the default tag (`"env"`) name for convertions.

```go
func ParseTagInto(tag string, c interface{})
```
Uses a custom tag passed as argument for convertions.

```go
func ParseTagIntoFunc(tag string, c interface{}, convFn ConvFunc)
```
Uses a custom conversion function when a tag-environment variable match is found, so to allow custom type management. See the [docs](https://godoc.org/github.com/ruggi/env) for more.

