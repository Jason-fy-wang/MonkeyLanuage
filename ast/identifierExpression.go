package ast

import "com.language/monkey/token"

// expression
type Identifier struct {
	Token token.Token
	Value string
}

func (it *Identifier) statementNode() {}

func (it *Identifier) TokenLiteral() string {
	return it.Token.Literal
}

func (it *Identifier) expressionNode() {
}

func (it *Identifier) String() string {
	return it.Value
}
