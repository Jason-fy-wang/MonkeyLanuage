package parser

import (
	"strings"
	"testing"

	"com.language/monkey/ast"
	"com.language/monkey/lexer"
	"com.language/monkey/token"
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

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world"`

	l := lexer.New(input)
	p := New(l)
	program := p.ParserProgram()

	CheckParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.StringLiteral)

	if !ok {
		t.Fatalf("expect StringLiteral, got %T", program.Statements[0].(*ast.ExpressionStatement).Expression)
	}

	if stmt.TokenLiteral() != "hello world" {
		t.Errorf("expect %s, got %s", input, stmt.TokenLiteral())
	}
}

func TestLetParse(t *testing.T) {
	tests := []struct {
		input            string
		expectIdentifier string
		expectVal        interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, itm := range tests {
		l := lexer.New(itm.input)
		p := New(l)
		program := p.ParserProgram()
		CheckParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("expect program statements length 1, got %d", len(program.Statements))
		}

		stmt := program.Statements[0].(*ast.LetStatement)
		if !testLetStatement(t, stmt, itm.expectIdentifier) {
			return
		}

		testingLiteralExpression(t, stmt.Value, itm.expectVal)
	}
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

	idt1 := program.Statements[0].(*ast.ReturnStatement).Value.(*ast.Identifier)
	if idt1.Value != "null" {
		t.Errorf("expect null, got %s", idt1.Value)
	}

	idt2 := program.Statements[1].(*ast.ReturnStatement).Value.(*ast.CallExpression)

	if idt2.Function.(*ast.Identifier).Value != "add" {
		t.Errorf("expect add identifier, got %s", idt2.Function)
	}

	if len(idt2.Arguments) != 2 {
		t.Errorf("expect arguments length 2, got %d", len(idt2.Arguments))
	}

	idt3 := program.Statements[2].(*ast.ReturnStatement).Value.(*ast.IntegerLiteral)

	if idt3.Value != 3 {
		t.Errorf("expect %d, got %d", 3, idt3.Value)
	}
}
