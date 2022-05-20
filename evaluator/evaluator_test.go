package evaluator

import (
	"github.com/idexter/monkey/lexer"
	"github.com/idexter/monkey/object"
	"github.com/idexter/monkey/parser"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !assert.True(t, ok) {
		return false
	}
	return assert.Equal(t, expected, result.Value)
}
