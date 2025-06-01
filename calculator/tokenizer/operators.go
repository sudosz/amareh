package tokenizer

import (
	"math"
	"strconv"

	"golang.org/x/exp/constraints"
)

var Operators = map[TokenType]Operator{
	PLUS:                  add,
	MINUS:                 subtract,
	MULTIPLY:              multiply,
	DIVIDE:                divide,
	MOD:                   modulo,
	CARET:                 pow,
	AMPERSAND:             bitwiseAnd,
	PIPE:                  bitwiseOr,
	EQUAL:                 equal,
	GREATER_THAN:          greaterThan,
	GREATER_THAN_OR_EQUAL: greaterThanOrEqual,
	LESS_THAN:             lessThan,
	LESS_THAN_OR_EQUAL:    lessThanOrEqual,
}

func token2Float64(t Token) float64 {
	return t.Value.(float64)
}

func number2Token[T constraints.Integer | constraints.Float](f T) Token {
	return Token{
		Type:     DECIMAL,
		rawValue: strconv.FormatFloat(float64(f), 'f', -1, 64),
		Value:    f,
	}
}

func add(a, b Token) (Token, error) {
	return number2Token(token2Float64(a) + token2Float64(b)), nil
}
func subtract(a, b Token) (Token, error) {
	return number2Token(token2Float64(a) - token2Float64(b)), nil
}
func multiply(a, b Token) (Token, error) {
	return number2Token(token2Float64(a) * token2Float64(b)), nil
}
func divide(a, b Token) (Token, error) {
	return number2Token(token2Float64(a) / token2Float64(b)), nil
}
func modulo(a, b Token) (Token, error) {
	return number2Token(math.Mod(token2Float64(a), token2Float64(b))), nil
}
func pow(a, b Token) (Token, error) {
	return number2Token(math.Pow(token2Float64(a), token2Float64(b))), nil
}
func bitwiseAnd(a, b Token) (Token, error) {
	t1 := int64(token2Float64(a))
	t2 := int64(token2Float64(b))
	if token2Float64(a)-float64(t1) > 0 || token2Float64(b)-float64(t2) > 0 {
		return number2Token(0), ErrInvalidDecimal
	}
	return number2Token(t1 & t2), nil
}
func bitwiseOr(a, b Token) (Token, error) {
	t1 := int64(token2Float64(a))
	t2 := int64(token2Float64(b))
	if token2Float64(a)-float64(t1) > 0 || token2Float64(b)-float64(t2) > 0 {
		return number2Token(0), ErrInvalidDecimal
	}
	return number2Token(t1 | t2), nil
}

func equal(a, b Token) (Token, error) {
	return Booleans[token2Float64(a) == token2Float64(b)], nil
}
func greaterThan(a, b Token) (Token, error) {
	return Booleans[token2Float64(a) > token2Float64(b)], nil
}
func greaterThanOrEqual(a, b Token) (Token, error) {
	return Booleans[token2Float64(a) >= token2Float64(b)], nil
}
func lessThan(a, b Token) (Token, error) {
	return Booleans[token2Float64(a) < token2Float64(b)], nil
}
func lessThanOrEqual(a, b Token) (Token, error) {
	return Booleans[token2Float64(a) <= token2Float64(b)], nil
}
