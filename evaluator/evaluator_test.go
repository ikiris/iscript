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
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)
		testIntegerObject(t, got, tt.expected)
	}
}

func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	env := object.NewEnvironment()
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("got err %v", err)
	}

	return Eval(prog, env)
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
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

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)
		i, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, got, int64(i))
		} else {
			testNullObj(t, got)
		}
	}
}

func testNullObj(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
		if (10 > 1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}`, 10},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)
		testIntegerObject(t, got, tt.want)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)

		errObj, ok := got.(*object.Error)
		if !ok {
			t.Errorf("no error returned. got=%T(%+v)", got, got)
			continue
		}

		if errObj.Message != tt.want {
			t.Errorf("wrong err message. got=`%q` want=`%q`", errObj.Message, tt.want)
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.want)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"

	got := testEval(t, input)
	fn, ok := got.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", got, got)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("func has incorrect params. got=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("paramater is not 'x'. got=%s", fn.Parameters[0])
	}

	wantBody := "(x + 2)"
	if fn.Body.String() != wantBody {
		t.Fatalf("body is not %q. got=%q", wantBody, fn.Body)
	}
}

func TestFuncApplication(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y }; add(5,5);", 10},
		{"let add = fn(x, y) { x + y }; add(5+5,add(5,5));", 20},
		{"fn(x, y) { x; }; (5);", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.want)
	}
}

func TestClosure(t *testing.T) {
	input := `
let newAdder = fn(x) {
	fn(y) {x + y};
};

let addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(t, input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	got := testEval(t, input)
	str, ok := got.(*object.String)
	if !ok {
		t.Fatalf("object is not string. got=%T (%+v)", got, got)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}
