package parser

import (
	"fmt"
	"iscript/ast"
	"iscript/lexer"
	"iscript/token"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestParserStruct(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *ast.Program
		wantErr  bool
	}{
		{
			"LetStatement",
			`
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`,
			&ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Token: token.Token{Type: "LET", Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: "IDENT", Literal: "x"},
							Value: "x",
						},
						Value: &ast.IntegerLiteral{Token: token.Token{Type: "INT", Literal: "5"}, Value: 5},
					},
					&ast.LetStatement{
						Token: token.Token{Type: "LET", Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: "IDENT", Literal: "y"},
							Value: "y",
						},
						Value: &ast.IntegerLiteral{Token: token.Token{Type: "INT", Literal: "10"}, Value: 10},
					},
					&ast.LetStatement{
						Token: token.Token{Type: "LET", Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: "IDENT", Literal: "foobar"},
							Value: "foobar",
						},
						Value: &ast.IntegerLiteral{Token: token.Token{Type: "INT", Literal: "838383"}, Value: 838383},
					},
				},
			},
			false,
		},
		{
			"ReturnStatement",
			`
			return 5;
			return 10;
			return 993322;
	`,
			&ast.Program{
				Statements: []ast.Statement{
					&ast.ReturnStatement{
						Token:       token.Token{Type: "RETURN", Literal: "return"},
						ReturnValue: &ast.IntegerLiteral{Token: token.Token{Type: "INT", Literal: "5"}, Value: 5},
					},
					&ast.ReturnStatement{
						Token:       token.Token{Type: "RETURN", Literal: "return"},
						ReturnValue: &ast.IntegerLiteral{Token: token.Token{Type: "INT", Literal: "10"}, Value: 10},
					},
					&ast.ReturnStatement{
						Token:       token.Token{Type: "RETURN", Literal: "return"},
						ReturnValue: &ast.IntegerLiteral{Token: token.Token{Type: "INT", Literal: "993322"}, Value: 993322},
					},
				},
			},
			false,
		},
		{
			"ExpressionStatement",
			`foobar;`,
			&ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: "IDENT", Literal: "foobar"},
						Expression: &ast.Identifier{
							Token: token.Token{Type: "IDENT", Literal: "foobar"},
							Value: "foobar",
						},
					},
				},
			},
			false,
		},
		{
			"integer literal",
			`5;`,
			&ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: "INT", Literal: "5"},
						Expression: &ast.IntegerLiteral{
							Token: token.Token{Type: "INT", Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			false,
		},
		{
			"prefix expressions",
			`!5;
			-15;`,
			&ast.Program{
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token: token.Token{Type: "!", Literal: "!"},
						Expression: &ast.PrefixExpression{
							Token:    token.Token{Type: "!", Literal: "!"},
							Operator: "!",
							Right: &ast.IntegerLiteral{
								Token: token.Token{
									Type:    "INT",
									Literal: "5",
								},
								Value: 5,
							},
						},
					},
					&ast.ExpressionStatement{
						Token: token.Token{Type: "-", Literal: "-"},
						Expression: &ast.PrefixExpression{
							Token:    token.Token{Type: "-", Literal: "-"},
							Operator: "-",
							Right: &ast.IntegerLiteral{
								Token: token.Token{
									Type:    "INT",
									Literal: "15",
								},
								Value: 15,
							},
						},
					},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		got, err := p.ParseProgram()
		if (err != nil) != tt.wantErr {
			t.Errorf("%s: got err %v when wantErr is %v", tt.name, err, tt.wantErr)
		}
		if diff := pretty.Compare(got, tt.expected); diff != "" {
			t.Errorf("%s: NextToken diff: (-got +want)\n%s", tt.name, diff)
		}
	}
}

func TestInfixParsing(t *testing.T) {
	tests := []struct {
		input    string
		leftVal  int64
		operator string
		rightVal int64
		wantErr  bool
	}{
		{"5 + 5;", 5, "+", 5, false},
		{"5 - 5;", 5, "-", 5, false},
		{"5 * 5;", 5, "*", 5, false},
		{"5 / 5;", 5, "/", 5, false},
		{"5 > 5;", 5, ">", 5, false},
		{"5 < 5;", 5, "<", 5, false},
		{"5 == 5;", 5, "==", 5, false},
		{"5 != 5;", 5, "!=", 5, false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program, err := p.ParseProgram()
		if err != nil {
			t.Errorf("%s: got err %v when wantErr is %v", tt.input, err, tt.wantErr)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("%s - program.Statements does not contain %d statements. got=%d", tt.input, 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("%s - program.Statements[0] is not ast.InfixExpression. got=%T", tt.input, stmt.Expression)
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("%s - exp is not ast.InfixExpression. got=%T", tt.input, stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
	}
}

func TestOpPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"!(true == true)", "(!(true == true))"},
		{"-(5+5)", "(-(5 + 5))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{"add(a * b[2], b[1], 2 * [1, 2][1])", "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program, err := p.ParseProgram()
		if err != nil {
			t.Fatalf("%s - failed to parse program: err: %v", tt.input, err)
		}

		got := program.String()
		if diff := pretty.Compare(got, tt.expected); diff != "" {
			t.Errorf("%s: NextToken diff: (-got +want)\n%s", tt.input, diff)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatalf("failed to parse program: err: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not %d statements. got=%d", 1, len(exp.Consequence.Statements))
	}

	c, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement0 is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, c.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements is not nil. got=%+v", exp.Alternative)
	}
}

func testIntLiteral(t *testing.T, l ast.Expression, v int64) bool {
	i, ok := l.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("l not *ast.IntegerLiteral. got=%T", l)
		return false
	}

	if i.Value != v {
		t.Errorf("i.Value not %d. got=%d", v, i.Value)
		return false
	}

	if i.TokenLiteral() != fmt.Sprintf("%d", v) {
		t.Errorf("i.TokenLiteral not %d. got=%s", v, i.TokenLiteral())
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, v string) bool {
	i, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Idenifier. got=%T", exp)
		return false
	}

	if i.Value != v {
		t.Errorf("i.Value not %s. got=%s", v, i.Value)
		return false
	}

	if i.TokenLiteral() != v {
		t.Errorf("i.TokenLiteral not %s. got=%s", v, i.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntLiteral(t, exp, int64(v))
	case int64:
		return testIntLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	return testLiteralExpression(t, opExp.Right, right)
}

func TestFuncLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatalf("failed to parse program: err: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Errorf("program is not %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}

	f, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(f.Parameters) != 2 {
		t.Fatalf("func paramters wrong, want %d, got=%d\n", 2, len(f.Parameters))
	}

	testLiteralExpression(t, f.Parameters[0], "x")
	testLiteralExpression(t, f.Parameters[1], "y")

	if len(f.Body.Statements) != 1 {
		t.Errorf("f.body.statements is not %d statements. got=%d", 1, len(f.Body.Statements))
	}

	bodyStmt, ok := f.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("f body statement is not ast.ExpressionStatement. got=%T", f.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCallFunc(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("failed to parse program: err: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Errorf("program is not %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong arg length. want %d got=%d", 3, len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("failed to parse program: err: %v", err)
	}

	if len(prog.Statements) != 1 {
		t.Errorf("program is not %d statements. got=%d", 1, len(prog.Statements))
	}

	stmt := prog.Statements[0].(*ast.ExpressionStatement)
	str, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.StringLiteral. got=%T", stmt.Expression)
	}

	if str.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", str.Value)
	}
}

func TestParseArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("failed to parse program: err: %v", err)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("stmt is not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len array not 3. got=%d", len(array.Elements))
	}

	testIntLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpression(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("failed to parse program: err: %v", err)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}

	idx, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp is not ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, idx.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, idx.Index, 1, "+", 1) {
		return
	}
}

func TestParseHashLiteralString(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("failed to parse program: err: %v", err)
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}

	hl, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	want := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hl.Pairs {
		lit, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key not stringliteral. got=%T", key)
		}

		expectedV := want[lit.String()]

		testIntLiteral(t, value, expectedV)
	}
}

func TestFuncLiteralWithName(t *testing.T) {
	input := `let myFunc = fn() {};`

	l := lexer.New(input)
	p := New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("failed to parse program: err: %v", err)
	}

	if len(prog.Statements) != 1 {
		t.Fatalf("prog.Body does not contain %d statements. got=%d\n", 1, len(prog.Statements))
	}

	stmt, ok := prog.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("prog.Statements[0] is not ast.LetStatement. got=%T", prog.Statements[0])
	}

	fn, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.FunctionLiteral. got=%T", stmt.Value)
	}

	if fn.Name != "myFunc" {
		t.Fatalf("function literal name wrong. want 'myFunc', got=%q", fn.Name)
	}
}
