package tokenizer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrUnexpectedCharacter = fmt.Errorf("unexpected character")
	ErrInvalidDecimal      = fmt.Errorf("invalid decimal")
	ErrInvalidExpession    = fmt.Errorf("invalid expression")
)

const (
	pi  = math.Pi
	phi = math.Phi
	e   = math.E
)

type Operator func(Token, Token) (Token, error)
type Function func(Token) Token

type Lexer struct {
	pos int
	exp []rune
}

func NewLexer(expression []rune) *Lexer {
	return &Lexer{
		pos: 0,
		exp: expression,
	}
}

func (l *Lexer) Lex() ([]Token, error) {
	tokens := make([]Token, 0)

	for l.pos < len(l.exp) {
		r := rune(l.exp[l.pos])
		fmt.Printf("Read: %c, IsDigit: %t, IsSpace: %t, IsOperator: %v \n", r, unicode.IsDigit(r), unicode.IsSpace(r), canOperator(r))

		if unicode.IsDigit(r) || r == '.' {
			token, err := l.lexDecimal()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		} else if unicode.IsSpace(r) {
			continue
		} else if op := canOperator(r); op != ILLEGAL {
			token, err := l.lexOperator(op)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		} else if token, cl := canConstant(l.exp[l.pos:]); token != Illegal {
			fmt.Printf("Constant: %v\n", token.Value)
			l.pos += cl
			tokens = append(tokens, token)
			continue
		} else {
			tokens = append(tokens, Illegal)
		}
		l.pos++
	}

	return tokens, nil
}

func canOperator(r rune) TokenType {
	if op, ok := operatorsTokenString[r]; ok {
		return op
	}
	return ILLEGAL
}

func canConstant(r []rune) (Token, int) {
	for c, t := range constantsTokenString {
		if strings.HasPrefix(string(r), c) {
			return t, len(c)
		}
	}
	return Illegal, 0
}

func (l *Lexer) lexConstant(c TokenType) (t Token, err error) {
	l.pos += len(Constants[c].rawValue)-1
	return Constants[c], nil
}

func (l *Lexer) lexOperator(op TokenType) (t Token, err error) {
	t.Type = op
	t.rawValue = op.String()
	t.Value = Operators[t.Type]
	if l.pos < len(l.exp)-1 {
		if t.Type == GREATER_THAN {
			if l.exp[l.pos+1] == '=' {
				t.Type = GREATER_THAN_OR_EQUAL
				t.rawValue += string(l.exp[l.pos+1])
				l.pos += 1
			} else {
				return t, nil
			}
		}
		if t.Type == LESS_THAN {
			if l.exp[l.pos+1] == '=' {
				t.Type = LESS_THAN_OR_EQUAL
				t.rawValue += string(l.exp[l.pos+1])
				l.pos += 1
			} else {
				return t, nil
			}
		}
	}

	return t, nil
}

func (l *Lexer) lexDecimal() (t Token, err error) {
	t.Type = DECIMAL

loop:
	for l.pos < len(l.exp) {
		r := rune(l.exp[l.pos])
		t.rawValue += string(r)
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
		case 'e':
			if t.rawValue == "e" {
				t.Value = math.E
				return t, nil
			}
		case '%':
			if t.rawValue == "%" {
				return t, fmt.Errorf("%w: %c", ErrUnexpectedCharacter, r)
			}
			if l.pos < len(l.exp)-1 {
				nextR := rune(l.exp[l.pos+1])
				if unicode.IsDigit(nextR) {
					t.rawValue = t.rawValue[:len(t.rawValue)-1]
					break loop
				}
			}
			val, err := strconv.ParseFloat(t.rawValue[:len(t.rawValue)-1], 64)
			if err != nil {
				return t, err
			}
			t.Value = val * 0.01
			return t, nil
		case ',':
			if len(t.rawValue) > 0 && t.rawValue[len(t.rawValue)-1] == ',' {
				return t, fmt.Errorf("%w: %c", ErrUnexpectedCharacter, r)
			}
		case ' ':
			t.rawValue = t.rawValue[:len(t.rawValue)-1]
			break loop
		default:
			t.rawValue = t.rawValue[:len(t.rawValue)-1]
			break loop
		}
		l.pos++
	}

	t.Value, err = strconv.ParseFloat(t.rawValue, 64)
	if err != nil {
		return t, err
	}
	l.pos--

	return t, nil
}

func Tokenize(expression []rune) ([]Token, error) {
	lexer := NewLexer(expression)
	return lexer.Lex()
}
