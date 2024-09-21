package ast

import (
	"bytes"
	"strings"

	"com.lanuage/monkey/token"
)

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatements
}

func (fl *FunctionLiteral) expressionNode() {

}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fl.Token.Literal)
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	//out.WriteString("{")
	out.WriteString(fl.Body.String())
	//out.WriteString("}")

	return out.String()
}
