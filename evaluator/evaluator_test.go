package evaluator

import (
	"iscript/lexer"
	"iscript/object"
	"iscript/parser"
	"testing"
)

func TestIntEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)
		testIntegerObject(t, got, tt.expected)
	}
}

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("got err %v", err)
	}

	return Eval(prog)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj is not integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("value not as expected. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBool(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)
		testBoolObj(t, got, tt.want)
	}
}

func testBoolObj(t *testing.T, obj object.Object, want bool) bool {
	got, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj is not boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if got.Value != want {
		t.Errorf("obj has wrong value. got=%t, want=%t", got, want)
		return false
	}
	return true
}

func TestBang(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)
		testBoolObj(t, got, tt.want)
	}
}
