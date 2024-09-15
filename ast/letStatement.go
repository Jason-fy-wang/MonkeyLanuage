package ast

import (
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

func (lt *LetStatement) TokenLiteral() string {
	return lt.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (it *Identifier) statementNode() {}

func (it *Identifier) TokenLiteral() string {
	return it.Token.Literal
}
