package ast

import (
	"bytes"

	"com.lanuage/monkey/token"
)

/*
let identifier = <expression>;
*/
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (lt *LetStatement) statementNode() {

}

func (lt *LetStatement) String() string {
	var buf = bytes.Buffer{}

	buf.WriteString(lt.TokenLiteral() + " ")
	buf.WriteString(lt.Name.String())
	buf.WriteString(" = ")
	if lt.Value != nil {
		buf.WriteString(lt.Value.String())
	}
	buf.WriteString(";")

	return buf.String()
}

func (lt *LetStatement) TokenLiteral() string {
	return lt.Token.Literal
}
