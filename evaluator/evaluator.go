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
	if node == nil {
		return nil
	}
	switch nod := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: nod.Value}

	case *ast.Boolean:
		return nativeBooltoToBooleanObject(nod.Value)

	case *ast.PrefixExpression:
		right := Eval(nod.Right)
		return evalPrefixExpression(nod.Operator, right)

	case *ast.InFixExpression:
		left := Eval(nod.Left)
		right := Eval(nod.Right)
		return evalInfixExpression(nod.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(nod)

	case *ast.BlockStatements:
		if nod.Statements != nil {
			return evalStatements(nod.Statements)
		}
		return nil
	case *ast.Program:
		return evalStatements(nod.Statements)

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

func evalInfixExpression(operator string, left, right object.Object) object.Object {

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBooltoToBooleanObject(left == right)
	case operator == "!=":
		return nativeBooltoToBooleanObject(left != right)
	default:
		return NULL
	}

}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	case "<":
		return nativeBooltoToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBooltoToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBooltoToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBooltoToBooleanObject(leftVal != rightVal)
	default:
		return NULL

	}
}

func evalIfExpression(node *ast.IfExpression) object.Object {

	condition := Eval(node.Confition)

	if isTruthy(condition) {
		return Eval(node.Consequence)
	} else if !isTruthy(condition) {
		if node.Alternative != nil {
			return Eval(node.Alternative)
		} else {
			return NULL
		}
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {

	switch obj {
	case NULL:
		return false
	case TRUE:
		return true

	case FALSE:
		return false
	default:
		return true
	}
}

func evalStatements(stmts []ast.Statement) object.Object {

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
