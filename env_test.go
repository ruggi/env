package env_test

import (
	"os"
	"testing"

	"github.com/ruggi/env"
	"github.com/stretchr/testify/assert"
)

type v struct {
	key string
	val string
}

func TestParse(t *testing.T) {
	type T struct {
		A string  `env:"A"`
		B int     `env:"B"`
		C bool    `env:"C"`
		D byte    `env:"D"` // unsupported
		E int8    `env:"E"`
		F int16   `env:"F"`
		G int32   `env:"G"`
		H int64   `env:"H"`
		I float32 `env:"I"`
		J float64 `env:"J"`
	}
	tests := []struct {
		vars     []v
		c        T
		expected T
	}{
		{
			c:        T{A: "keep"},
			expected: T{A: "keep"},
		},
		{
			vars: []v{
				v{key: "A", val: "test!"},
			},
			c:        T{A: "replaceme"},
			expected: T{A: "test!"},
		},
		{
			vars: []v{
				v{key: "B", val: "wrong"},
			},
			c:        T{B: 42},
			expected: T{B: 42},
		},
		{
			vars: []v{
				v{key: "B", val: "9000"},
			},
			c:        T{B: 42},
			expected: T{B: 9000},
		},
		{
			vars: []v{
				v{key: "A", val: "test"},
				v{key: "B", val: "9000"},
				v{key: "C", val: "1"},
				v{key: "D", val: "0x1"},
			},
			c:        T{},
			expected: T{A: "test", B: 9000, C: true},
		},
		{
			vars: []v{
				v{key: "E", val: "2"},
				v{key: "F", val: "3"},
				v{key: "G", val: "4"},
				v{key: "H", val: "5"},
				v{key: "I", val: "1.123"},
				v{key: "J", val: "1.123456"},
			},
			c:        T{},
			expected: T{E: 2, F: 3, G: 4, H: 5, I: 1.123, J: 1.123456},
		},
	}
	for _, tt := range tests {
		for _, vv := range tt.vars {
			os.Setenv(vv.key, vv.val)
		}
		env.ParseInto(&tt.c)
		assert.Equal(t, tt.expected, tt.c)
		for _, vv := range tt.vars {
			os.Setenv(vv.key, "")
		}
	}
}
