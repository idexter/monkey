package parser

import (
	"fmt"
	"testing"

	"github.com/idexter/monkey/ast"
	"github.com/idexter/monkey/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		require.NotNil(t, program, "ParseProgram() returned nil")
		require.Len(t, program.Statements, 1, "program.Statements does not contain 1 statements.")

		stmt := program.Statements[0]
		testLetStatement(t, stmt, tt.expectedIdentifier)

		val := stmt.(*ast.LetStatement).Value
		testLiteralExpression(t, val, tt.expectedValue)
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		require.Len(t, program.Statements, 1, "program has not enough statements.")
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement.")

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		require.True(t, ok, "stmt is not ast.PrefixExpression.")

		assert.Equal(t, tt.operator, exp.Operator, "exp.Operator is not '%s'", tt.operator)
		testLiteralExpression(t, exp.Right, tt.integerValue)
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	integ, ok := il.(*ast.IntegerLiteral)
	require.True(t, ok, "il not *ast.IntegerLiteral.")
	require.Equal(t, value, integ.Value, "integ.Value not %d.", value)
	require.Equal(t, fmt.Sprintf("%d", value), integ.TokenLiteral(), "integ.TokenLiteral not %d.", value)
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
		require.Equal(t, tt.operator, exp.Operator, "exp.Operator is not '%s'", tt.operator)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		assert.Equal(t, tt.expected, actual)
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) {
	ident, ok := exp.(*ast.Identifier)
	require.True(t, ok, "exp not *ast.Identifier.")
	require.Equal(t, value, ident.Value, "ident.Value not %s.", value)
	require.Equal(t, value, ident.TokenLiteral(), "ident.TokenLiteral not %s.", value)
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	switch v := expected.(type) {
	case int:
		testIntegerLiteral(t, exp, int64(v))
		return
	case int64:
		testIntegerLiteral(t, exp, v)
		return
	case string:
		testIdentifier(t, exp, v)
		return
	case bool:
		testBooleanLiteral(t, exp, v)
		return
	}
	t.Errorf("type of exp not handled. got=%T", exp)
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) {
	opExp, ok := exp.(*ast.InfixExpression)
	require.True(t, ok, "exp is not ast.InfixExpression.")

	testLiteralExpression(t, opExp.Left, left)
	require.Equal(t, operator, opExp.Operator, "exp.Operator is not '%s'", operator)
	testLiteralExpression(t, opExp.Right, right)
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	bo, ok := exp.(*ast.Boolean)
	require.True(t, ok, "exp not *ast.Boolean.")
	require.Equal(t, value, bo.Value, "bo.Value not %t.", value)
	require.Equal(t, fmt.Sprintf("%t", value), bo.TokenLiteral(), "bo.TokenLiteral not %t.", value)
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		require.Len(t, program.Statements, 1, "program has not enough statements.")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement.")

		boolean, ok := stmt.Expression.(*ast.Boolean)
		require.True(t, ok, "exp not *ast.Boolean.")

		require.Equal(t, tt.expectedBoolean, boolean.Value, "boolean.Value not %t.", tt.expectedBoolean)
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Len(t, program.Statements, 1, "program.Statements does not contain %d statements.", 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement.")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok, "stmt.Expression is not ast.IfExpression.")

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	require.Len(t, exp.Consequence.Statements, 1, "consequence is not 1 statements.")

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "Statements[0] is not ast.ExpressionStatement.")

	testIdentifier(t, consequence.Expression, "x")

	require.Nil(t, exp.Alternative, "exp.Alternative.Statements was not nil.")
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Len(t, program.Statements, 1, "program.Statements does not contain %d statements.", 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement.")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	require.True(t, ok, "stmt.Expression is not ast.IfExpression.")

	testInfixExpression(t, exp.Condition, "x", "<", "y")

	require.Len(t, exp.Consequence.Statements, 1, "consequence is not 1 statements.")

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "Statements[0] is not ast.ExpressionStatement.")

	testIdentifier(t, consequence.Expression, "x")

	require.Len(t, exp.Alternative.Statements, 1, "exp.Alternative.Statements does not contain 1 statements.")

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "Statements[0] is not ast.ExpressionStatement.")

	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	require.Len(t, program.Statements, 1, "program.Statements does not contain %d statements.", 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "program.Statements[0] is not ast.ExpressionStatement.")

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	require.True(t, ok, "stmt.Expression is not ast.FunctionLiteral.")

	require.Len(t, function.Parameters, 2, "function literal parameters wrong. want 2.")

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	require.Len(t, function.Body.Statements, 1, "function.Body.Statements has not 1 statements.")

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	require.True(t, ok, "function body stmt is not ast.ExpressionStatement.")

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	testIdentifier(t, exp.Function, "add")

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	testIdentifier(t, indexExp.Left, "myArray")
	testInfixExpression(t, indexExp.Index, 1, "+", 1)
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}
