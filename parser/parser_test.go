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
					},
					&ast.LetStatement{
						Token: token.Token{Type: "LET", Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: "IDENT", Literal: "y"},
							Value: "y",
						},
					},
					&ast.LetStatement{
						Token: token.Token{Type: "LET", Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: "IDENT", Literal: "foobar"},
							Value: "foobar",
						},
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
						Token: token.Token{Type: "RETURN", Literal: "return"},
					},
					&ast.ReturnStatement{
						Token: token.Token{Type: "RETURN", Literal: "return"},
					},
					&ast.ReturnStatement{
						Token: token.Token{Type: "RETURN", Literal: "return"},
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
