package lexer

import (
	"testing"

	"com.lanuage/monkey/token"
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
