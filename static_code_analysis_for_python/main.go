package main

type TokenType int

// Python Lexer

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

func (pl *PythonLexer) advance() {
	pl.position++
	pl.column++
}

func (pl *PythonLexer) skipWhiteSpace() {
	for pl.position < len(pl.input) && isWhiteSpace(pl.current()) {
		pl.advance()
	}
}

func (pl *PythonLexer) nextToken() Token {
	// TODO: write the next token
	return Token{}
}

// helper functions

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch == 'A' && ch == 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= 0 && ch <= '9'
}
