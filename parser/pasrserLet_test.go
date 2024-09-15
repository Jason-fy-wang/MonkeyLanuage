package parser

import (
	"strings"
	"testing"

	"com.lanuage/monkey/ast"
	"com.lanuage/monkey/lexer"
	"com.lanuage/monkey/token"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y= 10;
		let foobat = 828384;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParserProgram()
	CheckParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements didn't get 3 statements .  got : %d", len(program.Statements))
	}

	tests := []struct {
		expectIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobat"},
	}

	for i, itm := range tests {
		pt := program.Statements[i]

		if !testLetStatement(t, pt, itm.expectIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {

	if stmt.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not let. got= %q", stmt.TokenLiteral())
		return false
	}

	letstmt, ok := stmt.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not let statement. got %T", stmt)
		return false
	}

	if letstmt.Name.Value != name {
		t.Errorf("letStatement.Name.Value not %s, get %s", name, letstmt.Name.Value)
		return false
	}

	if letstmt.Name.TokenLiteral() != name {
		t.Errorf("letstmt.Name.TOkenLiteral() not %s, got %s", name, letstmt.Name.TokenLiteral())

		return false
	}

	return true
}

func CheckParserErrors(t *testing.T, p *Parser) {
	if len(p.errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.errors))

	errmsg := strings.Join(p.errors, "\n")
	t.Errorf("get error msgs %s", errmsg)
	t.FailNow()

}

func TestReturnStatements(t *testing.T) {
	inputs := `

	return null;
	return add(5,5);
	return 3;
	`

	lex := lexer.New(inputs)
	p := New(lex)

	program := p.ParserProgram()

	if program == nil {
		t.Fatal("parset returnStatement return nil.")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("parse returnStatment not 3;  get %d", len(program.Statements))
	}

	for _, itm := range program.Statements {
		rsstmt, ok := itm.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("parser didn't parse returnStatement.")
		}
		if rsstmt.Token.Type != token.RETURN {
			t.Fatalf("tokenType not RETURN. got %s", rsstmt.Token.Type)
		}

		if rsstmt.Token.Literal != "return" {
			t.Fatalf("tokenLiteral not return. got %s", rsstmt.Token.Literal)
		}
	}
}
