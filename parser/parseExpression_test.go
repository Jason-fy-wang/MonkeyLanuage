package parser

import (
	"fmt"
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp note *ast.Identifier, got %T", exp)

		return false
	}

	if ident.Value != value {
		t.Errorf("ident Value not %s,  got %s", value, ident.Value)

		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s, got %s", value, ident.TokenLiteral())

		return false
	}
	return true
}

func testIntegerLiterals(t *testing.T, il ast.Expression, val int64) bool {

	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not ast.IntegerLiteral, got %T", il)

		return false
	}

	if integ.Value != val {
		t.Errorf("intge.Value not %d, got %d", val, integ.Value)
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("integ.TokenLiteral not %d, got %s", val, integ.TokenLiteral())
		return false
	}

	return true
}

func testingLiteralExpression(t *testing.T, exp ast.Expression, expect interface{}) bool {

	switch v := expect.(type) {
	case int:
		return testIntegerLiterals(t, exp, int64(v))
	case int64:
		return testIntegerLiterals(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}

	t.Errorf("type of exp not handled. got %T", expect)

	return false
}

func testInfixExression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InFixExpression)

	if !ok {
		t.Errorf("exp is not ast.InfixExpression.  got %T(%s)", exp, exp)
		return false
	}

	if !testingLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %s, got =%q", operator, opExp.Operator)
		return false
	}

	if !testingLiteralExpression(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestBoolExpression(t *testing.T) {

	inputs := `
		true;
		false;
	`

	l := lexer.New(inputs)
	p := New(l)

	program := p.ParserProgram()

	if len(program.Statements) != 2 {
		t.Errorf("parse statements length not 2, got %d", len(program.Statements))
	}

	be, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("p.statement[0] not expressionStatement, got %T", program.Statements[0])
	}

	boolExp, ok := be.Expression.(*ast.Boolean)
	if !ok {
		t.Errorf("expressionStatement.expression not ast.Boolean. got %T", be.Expression)
	}

	if boolExp.Token.Type != token.TRUE {
		t.Errorf("expect token.TRUE, got: %v", boolExp.Token.Type)
	}

	if boolExp.Value != true {
		t.Errorf("expect true.  got %v", boolExp.Value)
	}
}
