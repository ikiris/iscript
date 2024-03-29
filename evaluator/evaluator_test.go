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
		{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
		{`{"name": "Monkey"}[fn(x) { x }];`, "unusable as hash key: FUNCTION"},
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

func TestStringConcat(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	got := testEval(t, input)
	str, ok := got.(*object.String)
	if !ok {
		t.Fatalf("object is not string. got=%T (%+v)", got, got)
	}

	if str.Value != "Hello World!" {
		t.Errorf("string has wrong value. got=%s", str.Value)
	}
}

func TestBuiltIn(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got=INTEGER"},
		{`len("one", "two")`, "wrong number of args: got=2, want=1"},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)

		switch want := tt.want.(type) {
		case int:
			testIntegerObject(t, got, int64(want))
		case string:
			errObj, ok := got.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", got, got)
				continue
			}
			if errObj.Message != want {
				t.Errorf("wrong error message. got=%q, want=%q", got, want)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	got := testEval(t, input)
	res, ok := got.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", got, got)
	}

	if len(res.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d", len(res.Elements))
	}

	testIntegerObject(t, res.Elements[0], 1)
	testIntegerObject(t, res.Elements[1], 4)
	testIntegerObject(t, res.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"[1, 2, 3][3]", nil},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)
		intg, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, got, int64(intg))
		} else {
			testNullObj(t, got)
		}
	}
}

func TestHashLiteral(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	got := testEval(t, input)
	result, ok := got.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return hash. got=%T (%+v)", got, got)
	}

	want := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(want) {
		t.Fatalf("pair count wrong. got=%d", len(result.Pairs))
	}

	for expectK, expectV := range want {
		pair, ok := result.Pairs[expectK]
		if !ok {
			t.Errorf("no pair for key")
		}

		testIntegerObject(t, pair.Value, expectV)
	}
}

func TestHashIndex(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		got := testEval(t, tt.input)
		intg, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, got, int64(intg))
		} else {
			testNullObj(t, got)
		}
	}
}

func TestMemoFib(t *testing.T) {
	tests := []struct {
		input string
		want  interface{}
	}{
		{
			input: `
			let cache = {};
			let memo = fn(f, x) {
				if (cache[x] != null) {
					return cache[x];
				};
				let c = f(x);
				updateHash(cache, x, c);
				return c;
			};
			let fib = fn(x) {
				if (x == 0) {
					return 0;
				};
				if (x == 1) {
					return 1;
				};
				memo(fib, x - 1) + memo(fib, x - 2);
			};
			memo(fib, 35);
			`,
			want: 9227465,
		},
	}
	for _, tt := range tests {
		got := testEval(t, tt.input)
		intg, ok := tt.want.(int)
		if ok {
			testIntegerObject(t, got, int64(intg))
		} else {
			testNullObj(t, got)
		}
	}
}
