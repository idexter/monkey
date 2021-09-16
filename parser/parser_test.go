package parser

import (
	"testing"

	"github.com/idexter/monkey/ast"
	"github.com/idexter/monkey/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLetStatements(t *testing.T) {
	input := `
   let x = 5;
   let y = 10;
   let foobar = 838383;
   `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.NotNil(t, program, "ParseProgram() returned nil")
	require.Len(t, program.Statements, 3, "program.Statements does not contain 3 statements.")

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(t, stmt, tt.expectedIdentifier)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	require.Equal(t, "let", s.TokenLiteral(), "s.TokenLiteral not 'let'")

	letStmt, ok := s.(*ast.LetStatement)
	require.True(t, ok, "s not *ast.LetStatement.")

	require.Equal(t, name, letStmt.Name.Value, "letStmt.Name.Value not '%s'", name)
	require.Equal(t, name, letStmt.Name.TokenLiteral(), "letStmt.Name.TokenLiteral() not '%s'", name)
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
   return 5;
   return 10;
   return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.NotNil(t, program, "ParseProgram() returned nil")
	require.Len(t, program.Statements, 3, "program.Statements does not contain 3 statements.")

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		assert.True(t, ok, "stmt not *ast.ReturnStatement.")
		assert.Equal(t, "return", returnStmt.TokenLiteral())
	}
}

func TestIdentifierStatement(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	require.NotNil(t, program, "ParseProgram() returned nil")

	require.Len(t, program.Statements, 1, "program has not enough statements.")
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement.")

	ident, ok := stmt.Expression.(*ast.Identifier)
	assert.True(t, ok, "exp not *ast.Identifier.")

	assert.Equal(t, "foobar", ident.Value, "ident.Value not %s", "foobar")
	assert.Equal(t, "foobar", ident.TokenLiteral(), "ident.TokenLiteral not %s", "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	require.NotNil(t, program, "ParseProgram() returned nil")

	require.Len(t, program.Statements, 1, "program has not enough statements.")
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement.")

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	assert.True(t, ok, "exp not *ast.Identifier.")

	assert.Equal(t, int64(5), ident.Value, "ident.Value not %s", "5")
	assert.Equal(t, "5", ident.TokenLiteral(), "ident.TokenLiteral not %s", "5")
}
