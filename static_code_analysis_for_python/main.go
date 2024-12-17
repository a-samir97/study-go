package main

import (
	"fmt"
	"os"
	"static_analysis/lexer"
)

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

// Python Analyzer

type Analyzer struct {
	lexer        *lexer.PythonLexer
	currentToken lexer.Token
	issues       []AnalyzerIssue
}

func NewAnalyzer(input string) *Analyzer {
	analyzer := Analyzer{
		lexer:  lexer.NewPythonLexer(input),
		issues: make([]AnalyzerIssue, 0),
	}
	analyzer.lexer.Advance()
	return &analyzer
}

func (a *Analyzer) advance() {
	a.currentToken = a.lexer.NextToken()
}

func (a *Analyzer) analyze() []AnalyzerIssue {
	for a.currentToken.Type != lexer.TOKEN_EOF {
		switch a.currentToken.Type {
		case lexer.TOKEN_DEF:
			a.analyzeFunctionDef()
		case lexer.TOKEN_IF:
			a.analyzeIfStatement(3)
		}
		a.advance()
	}
	return a.issues
}

func (a *Analyzer) analyzeFunctionDef() {
	// skip the 'def' keyword
	a.advance()

	if a.currentToken.Type != lexer.TOKEN_IDENTIFER {
		return
	}
	functionName := a.currentToken.Value

	// Count Parameters
	paramCount := 0

	for a.currentToken.Type != lexer.TOKEN_COLON && a.currentToken.Type != lexer.TOKEN_EOF {
		if a.currentToken.Type == lexer.TOKEN_IDENTIFER {
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
