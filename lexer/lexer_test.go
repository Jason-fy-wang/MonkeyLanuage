package lexer

import (
	"testing"

	"com.language/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		expectType    token.TokenType
		expectLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for _, itm := range tests {
		token := l.NextToken()

		if token.Literal != itm.expectLiteral && token.Type != itm.expectType {
			t.Errorf("parsed error. expect type: %v, real type: %s, expectLiteral: %s,  real literal: %s", itm.expectType, token.Type, itm.expectLiteral, token.Literal)
		}
	}
}

func TestLiteral(t *testing.T) {
	input := `let five=5;
				let ten=10;
				let add=fn(x,y){
					x+y
				};
				let result=add(five,ten);
				!-/*5
				5 < 10 > 5

				if (5 < 10) {
					return true;
				}else {
					return false;
				}

				10 == 10
				10 != 9
				10 <= 11
				10 >= 9
				"foobar"
				"foo bar"
				[1,2];
				{"foo":"bar"};
				`

	tests := []struct {
		expectType    token.TokenType
		expectLeteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.GREAT, ">"},
		{token.INT, "5"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LESS, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		// binocular operation (双目)
		{token.INT, "10"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.INT, "10"},
		{token.NOTEQUAL, "!="},
		{token.INT, "9"},
		{token.INT, "10"},
		{token.LEQ, "<="},
		{token.INT, "11"},
		{token.INT, "10"},
		{token.GEQ, ">="},
		{token.INT, "9"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		//{"foo":"bar"};
		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)

	for _, itm := range tests {
		token := l.NextToken()

		if itm.expectType != token.Type {
			t.Errorf("expect type: %s, real type: %v", itm.expectType, token.Type)
		}

		if itm.expectLeteral != token.Literal {
			t.Errorf("expect literal: %s, real literal: %v", itm.expectLeteral, token.Literal)
		}
	}
}

func TestDoubleOperation(t *testing.T) {
	input := `
	10 == 10
	10 != 9
	10 <= 11
	10 >= 9
	`

	tests := []struct {
		expectType    token.TokenType
		expectLeteral string
	}{
		// binocular operation (双目)
		{token.IDENT, "10"},
		{token.EQUAL, "=="},
		{token.IDENT, "10"},
		{token.IDENT, "10"},
		{token.NOTEQUAL, "!="},
		{token.IDENT, "9"},
		{token.IDENT, "10"},
		{token.LEQ, "<="},
		{token.IDENT, "11"},
		{token.IDENT, "10"},
		{token.GEQ, ">="},
		{token.IDENT, "9"},
	}

	l := New(input)

	for _, itm := range tests {
		token := l.NextToken()

		if itm.expectType != token.Type {
			t.Errorf("expect type: %s, real type: %v", itm.expectType, token.Type)
		}

		if itm.expectLeteral != token.Literal {
			t.Errorf("expect literal: %s, real literal: %v", itm.expectLeteral, token.Literal)
		}
	}
}
