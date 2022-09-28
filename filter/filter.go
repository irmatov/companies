// package filter transform a simple stack based filter language
// into an expression that cab be used as WHERE clause in a SQL query.
//
// Example filter:
//
//	name,"Apple",=
//	price,100,<
//	name,"Apple",=,price,100,<,or
package filter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Filter contains expression suitable for usage in WHERE clause of an SQL query. Arguments
// contain actual arguments for that clause that are used by the expression.
type Filter struct {
	Expr      string
	Arguments []interface{}
}

type stackValue struct {
	expr   string
	simple bool
}

func parseString(s string) (string, int) {
	if !strings.HasPrefix(s, `"`) {
		panic("string must start with a quote")
	}
	var slash bool
	runes := make([]rune, 0, len(s))
	for i, r := range s[1:] {
		switch r {
		case '\\':
			if slash {
				runes = append(runes, r)
				slash = false
			} else {
				slash = true
			}
		case '"':
			if slash {
				runes = append(runes, r)
				slash = false
			} else {
				return string(runes), i + 2
			}
		default:
			runes = append(runes, r)
			slash = false
		}
	}
	if slash {
		panic("unterminated slash escape sequence")
	}
	panic("string must end with a quote")
}

func parseWord(s string) (string, int) {
	if len(s) == 0 {
		panic("word must start with a letter")
	}
	for i, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			continue
		}
		return s[:i], i
	}
	return s, len(s)
}

func parseNumber(s string) (interface{}, int) {
	pos := strings.Index(s, ",")
	if pos == -1 {
		pos = len(s)
	}
	n, err := strconv.ParseInt(s[:pos], 10, 64)
	if err == nil {
		return int(n), pos
	}
	f, err := strconv.ParseFloat(s[:pos], 64)
	if err != nil {
		panic("invalid number")
	}
	return f, pos
}

type tokenType int

const (
	tokComma tokenType = iota
	tokIdentifier
	tokLiteral
	tokOperator
)

type token struct {
	kind  tokenType
	value interface{}
}

func tokens(expr string, ch chan<- token) {
	for pos := 0; pos < len(expr); {
		switch {
		case expr[pos] == ',':
			pos++
		case expr[pos] == '"':
			s, delta := parseString(expr[pos:])
			pos += delta
			ch <- token{tokLiteral, s}
		case expr[pos] >= 'a' && expr[pos] <= 'z':
			fallthrough
		case expr[pos] >= 'A' && expr[pos] <= 'Z':
			s, delta := parseWord(expr[pos:])
			pos += delta
			if s == "and" || s == "or" {
				ch <- token{tokOperator, s}
			} else {
				ch <- token{tokIdentifier, s}
			}
		case expr[pos] == '=':
			ch <- token{tokOperator, "="}
			pos++
		case expr[pos] == '<' || expr[pos] == '>':
			if pos+1 < len(expr) && expr[pos+1] == '=' {
				ch <- token{tokOperator, expr[pos : pos+2]}
				pos += 2
			} else {
				ch <- token{tokOperator, expr[pos : pos+1]}
				pos++
			}
		case (expr[pos] == '+' || expr[pos] == '-') && (pos+1 >= len(expr) || expr[pos+1] == ','):
			ch <- token{tokOperator, expr[pos : pos+1]}
			pos++
		case expr[pos] == '+' || expr[pos] == '-' || (expr[pos] >= '0' && expr[pos] <= '9'):
			n, delta := parseNumber(expr[pos:])
			pos += delta
			ch <- token{tokLiteral, n}
		default:
			panic(fmt.Sprintf("unexpected character: %q", expr[pos]))
		}
	}
}

func bracketedExpr(v stackValue) string {
	if v.simple {
		return v.expr
	}
	return fmt.Sprintf("(%s)", v.expr)
}

// Execute transforms given stack-based expression into a Filter.
// Field references are allowed only from the provided list of fields, otherwise
// an error is returned.
func Execute(fields []string, expr string) (Filter, error) {
	ch := make(chan token)
	var tokensError string
	go func() {
		defer func() {
			if err := recover(); err != nil {
				tokensError = err.(string)
			}
			close(ch)
		}()
		tokens(expr, ch)
	}()

	var stack []stackValue
	var args []interface{}
	knownFields := make(map[string]struct{})
	for _, field := range fields {
		knownFields[field] = struct{}{}
	}
	for t := range ch {
		switch t.kind {
		case tokComma:
		case tokLiteral:
			stack = append(stack, stackValue{fmt.Sprintf("$%d", len(args)+1), true})
			args = append(args, t.value)
		case tokIdentifier:
			if _, ok := knownFields[t.value.(string)]; !ok {
				return Filter{}, fmt.Errorf("unknown field: %q", t.value.(string))
			}
			stack = append(stack, stackValue{t.value.(string), true})
		case tokOperator:
			op := t.value.(string)
			switch op {
			case "=", "<", "<=", ">", ">=", "-", "+", "and", "or":
				if len(stack) < 2 {
					return Filter{}, fmt.Errorf("not enough arguments for operator %q", t.value)
				}
				left := stack[len(stack)-2]
				right := stack[len(stack)-1]
				expr := fmt.Sprintf("%s %s %s", bracketedExpr(left), op, bracketedExpr(right))
				stack[len(stack)-2] = stackValue{expr, false}
				stack = stack[0 : len(stack)-1 : cap(stack)]
			default:
				return Filter{}, fmt.Errorf("unknown operator %q", t.value)
			}
		}
	}
	if tokensError != "" {
		return Filter{}, errors.New(tokensError)
	}
	if len(stack) == 0 {
		return Filter{}, nil
	}
	if len(stack) != 1 {
		return Filter{}, fmt.Errorf("stack has %d elements, expected 1", len(stack))
	}
	return Filter{Expr: stack[0].expr, Arguments: args}, nil
}
