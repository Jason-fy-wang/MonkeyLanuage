package ast

import (
	"bytes"

	"com.language/monkey/token"
)

type IfExpression struct {
	Token       token.Token
	Confition   Expression
	Consequence *BlockStatements
	Alternative *BlockStatements
}

func (ife *IfExpression) expressionNode() {}

func (ife *IfExpression) TokenLiteral() string {
	return ife.Token.Literal
}

func (ife *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ife.Confition.String())
	out.WriteString(" ")
	out.WriteString(ife.Consequence.String())

	if ife.Alternative != nil {
		out.WriteString(ife.Alternative.String())
	}

	return out.String()
}

type BlockStatements struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatements) expressionNode() {}

func (bs *BlockStatements) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatements) String() string {
	var out bytes.Buffer

	for _, itm := range bs.Statements {
		out.WriteString(itm.String())
	}

	return out.String()
}
