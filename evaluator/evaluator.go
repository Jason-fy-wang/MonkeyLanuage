package evaluator

import (
	"fmt"
	"os"

	"com.lanuage/monkey/ast"
	"com.lanuage/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch nod := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: nod.Value}

	case *ast.Boolean:
		return nativeBooltoToBooleanObject(nod.Value)

	case *ast.Program:
		return evelStatements(nod.Statements)

	case *ast.ExpressionStatement:
		return Eval(nod.Expression)

	}
	fmt.Fprintf(os.Stderr, "invalid expression. %v", node)
	return nil
}

func evelStatements(stmts []ast.Statement) object.Object {

	var obj object.Object

	for _, itm := range stmts {
		obj = Eval(itm)
	}

	return obj
}

func nativeBooltoToBooleanObject(input bool) object.Object {

	if input {
		return TRUE
	}
	return FALSE
}
