package ast

import (
	"bytes"
	"strings"

	"com.language/monkey/token"
)

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, ele := range al.Elements {
		elements = append(elements, ele.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
