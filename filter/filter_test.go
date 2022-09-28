package filter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		fields     []string
		expression string
		want       Filter
		wantErr    error
	}{
		{
			"empty expression",
			nil,
			"",
			Filter{"", nil},
			nil,
		},
		{
			"empty expression but with fields",
			[]string{"first", "second"},
			"",
			Filter{"", nil},
			nil,
		},
		{
			"simple expression",
			[]string{"first"},
			`first,"value",=`,
			Filter{"first = $1", []interface{}{"value"}},
			nil,
		},
		{
			"simple expression with unknown field",
			[]string{},
			`first,"value",=`,
			Filter{},
			errors.New(`unknown field: "first"`),
		},
		{
			"simple expression with integer",
			[]string{"first"},
			`first,10,=`,
			Filter{"first = $1", []interface{}{10}},
			nil,
		},
		{
			"numbers with signs",
			[]string{},
			`-10,+10,=`,
			Filter{"$1 = $2", []interface{}{-10, 10}},
			nil,
		},
		{
			"simple expression with float",
			[]string{"first"},
			`first,10.5,1,-,=`,
			Filter{"first = ($1 - $2)", []interface{}{float64(10.5), 1}},
			nil,
		},
		{
			"complex expression",
			[]string{"first", "second", "third"},
			`first,10,<,second,"value",=,or,third,20,>=,and`,
			Filter{`((first < $1) or (second = $2)) and (third >= $3)`, []interface{}{10, "value", 20}},
			nil,
		},
		{
			"unexpected character",
			[]string{"first", "second"},
			`first second`,
			Filter{},
			errors.New("unexpected character: ' '"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Execute(tt.fields, tt.expression)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_parseString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
		skip int
		err  string
	}{
		{
			"empty string",
			"",
			"",
			0,
			"string must start with a quote",
		},
		{
			"string must start with a quote",
			"hello",
			"",
			0,
			"string must start with a quote",
		},
		{
			"unterminated string",
			`"hello`,
			"",
			0,
			"string must end with a quote",
		},
		{
			"simple string",
			`"hello"`,
			"hello",
			7,
			"",
		},
		{
			"string with a quotes inside",
			`"h\el\"l\\o"`,
			`hel"l\o`,
			12,
			"",
		},
		{
			"unterminated slash escape sequence",
			`"hello\`,
			"",
			0,
			"unterminated slash escape sequence",
		},
		{
			"simple string with leftover data",
			`"hello",something`,
			"hello",
			7,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var gotSkip int
			var panicMsg string
			func() {
				defer func() {
					if p := recover(); p != nil {
						panicMsg = p.(string)
					}
				}()
				got, gotSkip = parseString(tt.s)
			}()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.skip, gotSkip)
			assert.Equal(t, tt.err, panicMsg)
		})
	}
}

func Test_parseWord(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
		skip int
		err  string
	}{
		{
			"empty string",
			"",
			"",
			0,
			"word must start with a letter",
		},
		{
			"simple",
			"simple",
			"simple",
			6,
			"",
		},
		{
			"simple with extra",
			"simple,",
			"simple",
			6,
			"",
		},
		{
			"simple with extra",
			"simple#,",
			"simple",
			6,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			var gotSkip int
			var panicMsg string
			func() {
				defer func() {
					if p := recover(); p != nil {
						panicMsg = p.(string)
					}
				}()
				got, gotSkip = parseWord(tt.s)
			}()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.skip, gotSkip)
			assert.Equal(t, tt.err, panicMsg)
		})
	}
}
