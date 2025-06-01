package math

import (
	"fmt"
	"strings"

	"github.com/sudosz/amareh/calculator/tokenizer"
)

func solve(expression []rune) (string, error) {
	tokens, err := tokenizer.Tokenize(expression)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type.IsIllegal() {
			return "", tokenizer.ErrInvalidExpession
		}
		if tokens[i].Type.IsOperator() {
			if i == 0 || i == len(tokens)-1 {
				return "", tokenizer.ErrInvalidExpession
			}
			left, right := tokens[i-1], tokens[i+1]
			if left.Type != tokenizer.DECIMAL || right.Type != tokenizer.DECIMAL && !(left.Type.IsConstant() || right.Type.IsConstant()) {
				return "", tokenizer.ErrInvalidExpession
			}
			val, err := tokens[i].Value.(tokenizer.Operator)(left, right)
			if err != nil {
				return "", err
			}
			tokens[i+1] = val
			tokens = tokens[i+1:]
			i -= 2
		}
	}

	if len(tokens) != 1 {
		return "", tokenizer.ErrInvalidExpession
	}

	return fmt.Sprintf("%v", tokens[0].Value), nil
}

var fixReplacer = strings.NewReplacer("**", "^", "×", "*", "÷", "/", "∧", "^")

func fixExpression(expression string) []rune {
	return []rune(fixReplacer.Replace(expression))
}

func Solve(expression string) (string, error) {
	return solve(fixExpression(expression))
}
