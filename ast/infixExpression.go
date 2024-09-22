package ast

import (
	"bytes"

	"com.language/monkey/token"
)

type InFixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InFixExpression) expressionNode() {}

func (ie *InFixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *InFixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
