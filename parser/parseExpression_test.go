package parser

import (
	"testing"

	"com.lanuage/monkey/ast"
	"com.lanuage/monkey/lexer"
	"com.lanuage/monkey/token"
)

func TestIdentifierExpression(t *testing.T) {
	inputs := "foobar;"

	lex := lexer.New(inputs)

	p := New(lex)

	program := p.ParserProgram()

	if program == nil {
		t.Fatal("identifier parse nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("parse statements length: %d", len(program.Statements))
	}

	estmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expect expressionStateent, but got %v", program.Statements[0])
	}

	if estmt.Token.Type != token.IDENT {
		t.Fatalf("expect IDENT, but got: %s", estmt.Token.Type)
	}

	espresion, ok := estmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("can't eonvert to identifier")
	}

	if espresion.Token.Literal != "foobar" {
		t.Fatalf("expect foobar, got: %s", espresion.Value)
	}
}

func TestIntegerLiteral(t *testing.T) {
	inputs := "5;"

	lex := lexer.New(inputs)
	p := New(lex)

	program := p.ParserProgram()
	CheckParserErrors(t, p)
	if program == nil {
		t.Fatalf("parse program nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("expect statement length = 1, but got %d", len(program.Statements))
	}

	estmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatal("convert to expressionStatement fail")
	}

	if estmt.Token.Type != token.INT {
		t.Fatalf("expect INT type, got  %s", estmt.Token.Type)
	}

	ilexpression, ok := estmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatal("convert to IntegerLiteral fail")
	}

	if ilexpression.Token.Literal != "5" {
		t.Fatalf("expect `5`, but get: %s", ilexpression.TokenLiteral())
	}

	if ilexpression.Value != 5 {
		t.Fatalf("expect 5, got %d", ilexpression.Value)
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		inputs   string
		operator string
		intVal   int64
	}{
		{"-5", "-", 5},
		{"!10", "!", 10},
	}

	for _, itm := range prefixTests {
		lex := lexer.New(itm.inputs)
		p := New(lex)

		program := p.ParserProgram()
		CheckParserErrors(t, p)

		if program == nil {
			t.Fatal("program nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("statements length got %d", len(program.Statements))
		}

		estmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatal("concert to expressionStatement fail")
		}

		if estmt.Token.Literal != itm.operator {
			t.Fatalf("expect %s, got %s", itm.operator, estmt.Token.Literal)
		}

		prefix, ok := estmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatal("convert to prefixExpression failed")
		}

		if prefix.Operator != itm.operator {
			t.Fatalf("expect %s , got %s", itm.operator, prefix.Operator)
		}

		intval, ok := prefix.Right.(*ast.IntegerLiteral)
		if !ok {
			t.Fatal("convert to intLiteral fail")
		}

		if intval.Value != itm.intVal {
			t.Fatalf("expect %d, got %d", itm.intVal, intval.Value)
		}
	}
}
