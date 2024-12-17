package lexer

import "static_analysis/helpers"

// Python Lexer
type TokenType int

const (
	TOKEN_DEF TokenType = iota
	TOKEN_IF
	TOKEN_IDENTIFER
	TOKEN_COLON
	TOKEN_IDENT
	TOKEN_DEDENT
	TOKEN_NEWLINE
	TOKEN_EOF
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

type PythonLexer struct {
	input      string
	position   int
	line       int
	column     int
	identStack []int
}

func NewPythonLexer(input string) *PythonLexer {
	return &PythonLexer{
		input:      input,
		position:   0,
		line:       1,
		column:     1,
		identStack: []int{0},
	}
}

func (pl *PythonLexer) current() byte {
	if pl.position >= len(pl.input) {
		return 0
	}
	return pl.input[pl.position]
}

func (pl *PythonLexer) Advance() {
	pl.position++
	pl.column++
}

func (pl *PythonLexer) skipWhiteSpace() {
	for pl.position < len(pl.input) && helpers.IsWhiteSpace(pl.current()) {
		pl.Advance()
	}
}

func (pl *PythonLexer) match(s string) bool {
	if pl.position+len(s) > len(pl.input) {
		return false
	}

	if pl.input[pl.position:pl.position+len(s)] == s {
		pl.position += len(s)
		pl.column += len(s)
		return true
	}

	return false
}

func (pl *PythonLexer) NextToken() Token {
	pl.skipWhiteSpace()

	if pl.position >= len(pl.input) {
		return Token{Type: TOKEN_EOF}
	}

	switch pl.current() {
	case 'd':
		if pl.match("def") {
			return Token{Type: TOKEN_DEF, Value: "def", Line: pl.line, Column: pl.column}
		}
	case 'i':
		if pl.match("if") {
			return Token{Type: TOKEN_IF, Value: "if", Line: pl.line, Column: pl.column}
		}
	case ':':
		pl.Advance()
		return Token{Type: TOKEN_COLON, Value: ":", Line: pl.line, Column: pl.column}
	case '\n':
		pl.Advance()
		pl.line++
		pl.column = 1
		return Token{Type: TOKEN_NEWLINE, Line: pl.line}
	}
	// Handle Identifiers

	if helpers.IsLetter(pl.current()) {
		start := pl.position

		for pl.position < len(pl.input) && (helpers.IsLetter(pl.current()) || helpers.IsDigit(pl.current())) {
			pl.Advance()
		}
		value := pl.input[start:pl.position]
		return Token{Type: TOKEN_IDENTIFER, Value: value, Line: pl.line, Column: pl.column}
	}

	pl.Advance()
	return Token{Type: TOKEN_IDENTIFER, Value: string(pl.input[pl.position-1]), Line: pl.line, Column: pl.column}
}
