package parser

import (
	"iscript/ast"
	"iscript/lexer"
	"iscript/token"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestParser(t *testing.T) {
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
