package tokenizer

import "math"

type TokenType int

func (t TokenType) String() string {
	return tokenTypeStrings[t]
}

func (t TokenType) IsOperator() bool {
	switch t {
	case PLUS, MINUS, MULTIPLY, DIVIDE, PARENTHESIS_OPEN, PARENTHESIS_CLOSE, COMMA, SEMICOLON, COLON, MOD, CARET, AMPERSAND, PIPE, EQUAL, GREATER_THAN, GREATER_THAN_OR_EQUAL, LESS_THAN, LESS_THAN_OR_EQUAL:
		return true
	}
	return false
}

func (t TokenType) IsLogicalOperator() bool {
	switch t {
	case AND, OR, NOT:
		return true
	}
	return false
}

func (t TokenType) IsParanthesis() bool {
	switch t {
	case PARENTHESIS_OPEN, PARENTHESIS_CLOSE:
		return true
	}
	return false
}

func (t TokenType) IsIllegal() bool {
	return t == ILLEGAL
}

func (t TokenType) IsNaN() bool {
	return t == NOT_A_NUMBER
}

func (t TokenType) IsConstant() bool {
	switch t {
	case PHI, PI, E, INFINITY:
		return true
	}
	return false
}

const (
	EOF TokenType = iota
	ILLEGAL

	// Types
	DECIMAL // 1234567890
	PERCENT // 123%
	BOOLEAN // true/false

	// Operators
	// -- LOGICAL OPERATORS --
	AND // &&
	OR  // ||
	NOT // !
	// -- ARITHMETIC OPERATORS --
	PLUS              // +
	MINUS             // -
	MULTIPLY          // *
	DIVIDE            // /
	PARENTHESIS_OPEN  // (
	PARENTHESIS_CLOSE // )
	COMMA             // ,
	SEMICOLON         // ;
	COLON             // :
	MOD               // %
	CARET             // ^
	AMPERSAND         // &
	PIPE              // |
	// -- COMPARISON OPERATORS --
	EQUAL                 // =
	GREATER_THAN          // >
	GREATER_THAN_OR_EQUAL // >=
	LESS_THAN             // <
	LESS_THAN_OR_EQUAL    // <=

	// Constants
	PHI          // φ
	PI           // π
	E            // e
	INFINITY     // ∞
	NOT_A_NUMBER // NaN

	// Functions
	SIN       // sin
	COS       // cos
	TAN       // tan
	COT       // cot
	SEC       // sec
	CSC       // csc
	COSEC     // cosec
	ABS       // abs
	SQRT      // sqrt
	CBRT      // cbrt
	LOG       // log
	LN        // ln
	EXP       // exp
	FACTORIAL // !
	LIMIT     // lim

	// -- MATH OPERATORS --
	SUM        // Σ
	PRODUCT    // Π
	INTEGRAL   // ∫
	DERIVATIVE // ∂
)

var tokenTypeStrings = map[TokenType]string{
	EOF:     "EOF",     //
	ILLEGAL: "ILLEGAL", //

	// Types
	DECIMAL: "DECIMAL", //
	BOOLEAN: "BOOLEAN", //

	// Operators
	// -- LOGICAL OPERATORS --
	AND: "AND",
	OR:  "OR",
	NOT: "NOT",
	// -- ARITHMETIC OPERATORS --
	PLUS:              "+", //
	MINUS:             "-", ///
	MULTIPLY:          "*", //
	DIVIDE:            "/", //
	PARENTHESIS_OPEN:  "(",
	PARENTHESIS_CLOSE: ")",
	COMMA:             ",", ///
	SEMICOLON:         ";",
	COLON:             ":",
	MOD:               "%", //
	CARET:             "^", //
	AMPERSAND:         "&", //
	PIPE:              "|", //
	// -- COMPARISON OPERATORS --
	EQUAL:                 "=",  ///
	GREATER_THAN:          ">",  ///
	GREATER_THAN_OR_EQUAL: ">=", ///
	LESS_THAN:             "<",  ///
	LESS_THAN_OR_EQUAL:    "<=", ///

	// Constants
	PHI:          "φ",   ///
	PI:           "π",   ///
	E:            "e",   //
	INFINITY:     "∞",   //
	NOT_A_NUMBER: "NaN", //

	// Functions
	SIN:       "sin",
	COS:       "cos",
	TAN:       "tan",
	COT:       "cot",
	SEC:       "sec",
	CSC:       "csc",
	COSEC:     "cosec",
	ABS:       "abs",
	SQRT:      "sqrt",
	CBRT:      "cbrt",
	LOG:       "log",
	LN:        "ln",
	EXP:       "exp",
	LIMIT:     "lim",
	FACTORIAL: "!",

	// -- MATH OPERATORS --
	SUM:        "Σ",
	PRODUCT:    "Π",
	INTEGRAL:   "∫",
	DERIVATIVE: "∂",
}

var operatorsTokenString = map[rune]TokenType{
	'+': PLUS,
	'-': MINUS,
	'*': MULTIPLY,
	'/': DIVIDE,
	'(': PARENTHESIS_OPEN,
	')': PARENTHESIS_CLOSE,
	',': COMMA,
	';': SEMICOLON,
	':': COLON,
	'%': MOD,
	'^': CARET,
	'&': AMPERSAND,
	'|': PIPE,
	'=': EQUAL,
	'>': GREATER_THAN,
	'<': LESS_THAN,
}

var constantsTokenString = map[string]Token{
	"φ":   Constants[PHI],
	"phi": Constants[PHI],
	"π":   Constants[PI],
	"pi":  Constants[PI],
	"e":   Constants[E],
	"E":   Constants[E],
	"∞":   Constants[INFINITY],
	"inf": Constants[INFINITY],
	"nan": Constants[NOT_A_NUMBER],
}

type Token struct {
	Type     TokenType
	rawValue string
	Value    any // Operator/Function/Constant/Decimal
}

var (
	Booleans = map[bool]Token{
		true:  {Type: BOOLEAN, Value: true},
		false: {Type: BOOLEAN, Value: false},
	}
	Illegal   = Token{Type: ILLEGAL}
	Constants = map[TokenType]Token{
		PHI: {Type: PHI, Value: math.Phi},
		PI: {Type: PI, Value: math.Pi},
		E: {Type: E, Value: math.E},
		INFINITY: {Type: INFINITY, Value: math.Inf(1)},
		NOT_A_NUMBER: {Type: NOT_A_NUMBER, Value: math.NaN()},
	}
)
