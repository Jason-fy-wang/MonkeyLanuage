package ast

import (
	"bytes"

	"com.language/monkey/token"
)

// return <expression> ;
type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var buf = bytes.Buffer{}

	buf.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		buf.WriteString(rs.Value.String())
	}
	buf.WriteString(";")

	return buf.String()
}
