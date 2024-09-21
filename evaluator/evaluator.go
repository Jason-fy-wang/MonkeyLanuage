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

	case *ast.PrefixExpression:
		right := Eval(nod.Right)
		return evalPrefixExpression(nod.Operator, right)

	case *ast.Program:
		return evelStatements(nod.Statements)

	case *ast.ExpressionStatement:
		return Eval(nod.Expression)

	}
	fmt.Fprintf(os.Stderr, "invalid expression. %v", node)
	return nil
}

func evalPrefixExpression(operator string, rightNode object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(rightNode)
	case "-":
		return evalMinusPrefixOperationExpression(rightNode)

	default:
		return NULL
	}
}

func evalBangOperatorExpression(node object.Object) object.Object {
	switch node {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperationExpression(node object.Object) object.Object {
	if node.Type() != object.INTEGER_OBJ {
		return nil
	}

	val := node.(*object.Integer).Value
	return &object.Integer{Value: -val}
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
