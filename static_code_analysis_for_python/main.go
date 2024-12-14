package main

import (
	"fmt"
	"os"
)

type TokenType int

type AnalyzerIssue struct {
	Message  string
	Location string
}

type PythonAnalyzer struct {
	issues          []AnalyzerIssue
	currentFunction string
}

func NewPythonAnalyzer() *PythonAnalyzer {
	return &PythonAnalyzer{
		issues: make([]AnalyzerIssue, 0),
	}
}

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

func (pl *PythonLexer) nextToken() Token {
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
		pl.advance()
		return Token{Type: TOKEN_COLON, Value: ":", Line: pl.line, Column: pl.column}
	case '\n':
		pl.advance()
		pl.line++
		pl.column = 1
		return Token{Type: TOKEN_NEWLINE, Line: pl.line}
	}
	// Handle Identifiers

	if isLetter(pl.current()) {
		start := pl.position

		for pl.position < len(pl.input) && (isLetter(pl.current()) || isDigit(pl.current())) {
			pl.advance()
		}
		value := pl.input[start:pl.position]
		return Token{Type: TOKEN_IDENTIFER, Value: value, Line: pl.line, Column: pl.column}
	}

	pl.advance()
	return Token{Type: TOKEN_IDENTIFER, Value: string(pl.input[pl.position-1]), Line: pl.line, Column: pl.column}
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

// Python Analyzer

type Analyzer struct {
	lexer        *PythonLexer
	currentToken Token
	issues       []AnalyzerIssue
}

func NewAnalyzer(input string) *Analyzer {
	analyzer := Analyzer{
		lexer:  NewPythonLexer(input),
		issues: make([]AnalyzerIssue, 0),
	}
	analyzer.lexer.advance()
	return &analyzer
}

func (a *Analyzer) advance() {
	a.currentToken = a.lexer.nextToken()
}

func (a *Analyzer) analyze() []AnalyzerIssue {
	for a.currentToken.Type != TOKEN_EOF {
		switch a.currentToken.Type {
		case TOKEN_DEF:
			a.analyzeFunctionDef()
		case TOKEN_IF:
			a.analyzeIfStatement()
		}
		a.advance()
	}
	return a.issues
}

func (a *Analyzer) analyzeFunctionDef() {
	// skip the 'def' keyword
	a.advance()

	if a.currentToken.Type != TOKEN_IDENTIFER {
		return
	}
	functionName := a.currentToken.Value

	// Count Parameters
	paramCount := 0

	for a.currentToken.Type != TOKEN_COLON && a.currentToken.Type != TOKEN_EOF {
		if a.currentToken.Type == TOKEN_IDENTIFER {
			paramCount++
		}
		a.advance()
	}

	if paramCount > 5 {
		a.issues = append(a.issues, AnalyzerIssue{
			Message:  fmt.Sprintf("Function '%s' has too many parameters (%d)", functionName, paramCount),
			Location: fmt.Sprintf("line %d", a.currentToken.Line),
		})
	}
}

func (a *Analyzer) analyzeIfStatement(depth int) {
	// need to check
	if depth > 3 {
		a.issues = append(a.issues, AnalyzerIssue{
			Message:  "Deep nesting deteced (depth > 3)",
			Location: fmt.Sprintf("line %d", a.currentToken.Line),
		})
	}
}

func AnalyzePythonFile(filePath string) ([]AnalyzerIssue, error) {
	content, err := os.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("Failed to read the file")
	}

	analyzer := NewAnalyzer(string(content))

	return analyzer.analyze(), nil
}

func main() {
	filepath := "example.py"

	issues, err := AnalyzePythonFile(filepath)

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	fmt.Printf("Analysis complete. Found %d issues \n", len(issues))
	for _, issue := range issues {
		fmt.Printf("%s: %s\n", issue.Location, issue.Message)
	}
}
