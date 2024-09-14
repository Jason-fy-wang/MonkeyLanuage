package lexer

import (
	"fmt"
	"os"

	"com.lanuage/monkey/token"
)

type Lexer struct {
	Input        string
	Position     int
	ReadPosition int
	ch           byte
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	// skip space
	l.skipWhiteSpace()
	switch l.ch {
	case '=':
		tok = NewToken(token.ASSIGN, "=")
	case '+':
		tok = NewToken(token.PLUS, "+")
	case ',':
		tok = NewToken(token.COMMA, ",")
	case ';':
		tok = NewToken(token.SEMICOLON, ";")
	case '(':
		tok = NewToken(token.LPAREN, "(")
	case ')':
		tok = NewToken(token.RPAREN, ")")
	case '{':
		tok = NewToken(token.LBRACE, "{")
	case '}':
		tok = NewToken(token.RBRACE, "}")
	case '[':
		tok = NewToken(token.LBRACKET, "[")
	case ']':
		tok = NewToken(token.RBRACKET, "]")
	case 0:
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	default:
		if IsLitter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LoopupIdentifier(tok.Literal)
			return tok
		} else if IsDigital(l.ch) {
			tok.Literal = l.readDigital()
			tok.Type = token.LoopupIdentifier(tok.Literal)
			return tok
		} else {
			tok = NewToken(token.ILLEGAL, string(l.ch))
			fmt.Fprintf(os.Stderr, "invalid char: %v", l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.Position
	for IsLitter(l.ch) {
		l.readChar()
	}
	return l.Input[position:l.Position]
}

func (l *Lexer) readDigital() string {
	position := l.Position
	for IsDigital(l.ch) {
		l.readChar()
	}
	return l.Input[position:l.Position]
}

func (l *Lexer) readChar() {
	if l.ReadPosition < len(l.Input) {
		l.ch = l.Input[l.ReadPosition]
	} else {
		l.ch = 0
	}
	l.Position = l.ReadPosition
	l.ReadPosition = l.Position + 1
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.ReadPosition < len(l.Input) {
		return l.Input[l.ReadPosition]
	}

	return 0
}

func IsLitter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || (ch == '_')
}

func IsDigital(ch byte) bool {
	return ('0' <= ch && ch <= '9')
}

func NewToken(tType token.TokenType, literal string) token.Token {
	tk := token.Token{
		Type:    tType,
		Literal: literal,
	}
	return tk
}

func New(input string) *Lexer {
	lex := &Lexer{
		Input: input,
	}
	lex.readChar()
	return lex
}
