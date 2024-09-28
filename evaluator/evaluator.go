package evaluator

import (
	"fmt"

	"com.language/monkey/ast"
	"com.language/monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environement) object.Object {
	if node == nil {
		return nil
	}
	switch nod := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: nod.Value}

	case *ast.Boolean:
		return nativeBooltoToBooleanObject(nod.Value)

	case *ast.PrefixExpression:
		right := Eval(nod.Right, env)
		if IsError(right) {
			return right
		}
		return evalPrefixExpression(nod.Operator, right)

	case *ast.InFixExpression:
		left := Eval(nod.Left, env)
		if IsError(left) {
			return left
		}
		right := Eval(nod.Right, env)
		if IsError(right) {
			return right
		}

		return evalInfixExpression(nod.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(nod, env)

	case *ast.BlockStatements:
		if nod.Statements != nil {
			return evalBlockStatements(nod.Statements, env)
		}
		return nil
	case *ast.ReturnStatement:
		val := Eval(nod.Value, env)
		if IsError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(nod.Value, env)
		if IsError(val) {
			return val
		}
		env.Set(nod.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(nod, env)

	case *ast.StringLiteral:
		return &object.String{Value: nod.Value}

	case *ast.FunctionLiteral:
		params := nod.Parameters
		body := nod.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(nod.Function, env)

		if IsError(function) {
			return function
		}

		args := evalExpressions(nod.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(nod.Elements, env)

		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(nod.Left, env)
		if IsError(left) {
			return left
		}

		index := Eval(nod.Index, env)
		if IsError(index) {
			return index
		}

		return evalIndexExpression(left, index)

	case *ast.Program:
		return evalProgram(nod.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(nod.Expression, env)

	}
	//fmt.Fprintf(os.Stderr, "invalid expression. %v", node)
	return nil
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:

		return evalArrayIndexExpression(left, index)
	default:
		return NewError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(left, index object.Object) object.Object {
	arrayObj := left.(*object.Array)

	idx := index.(*object.Integer).Value

	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObj.Elements[idx]
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch function := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return function.Fn(args...)
	default:
		return NewError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environement {
	env := object.NewEnclosedEnvironment(fn.Env)

	for idx, param := range fn.Parameters {
		env.Set(param.Value, args[idx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Value
	}
	return obj
}

func evalExpressions(exps []ast.Expression, env *object.Environement) []object.Object {

	var result []object.Object

	for _, exp := range exps {
		obj := Eval(exp, env)
		if IsError(obj) {
			return []object.Object{obj}
		}

		result = append(result, obj)
	}
	return result
}

func evalPrefixExpression(operator string, rightNode object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(rightNode)
	case "-":
		return evalMinusPrefixOperationExpression(rightNode)

	default:
		return NewError("unknow operator:%s%s", operator, rightNode.Type())
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
		return NewError("unknow operator:-%s", node.Type())
	}

	val := node.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringinfixExpression(operator, left, right)
	case operator == "==":
		return nativeBooltoToBooleanObject(left == right)
	case operator == "!=":
		return nativeBooltoToBooleanObject(left != right)
	case left.Type() != right.Type():
		return NewError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return NewError("unknow operator: %s %s %s", left.Type(), operator, right.Type())
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
		return NewError("unknow operator: %s %s %s", left.Type(), operator, right.Type())

	}
}

func evalStringinfixExpression(operator string, left, right object.Object) object.Object {

	if operator != "+" {
		return NewError("unknow operator:%s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	return &object.String{Value: leftVal + rightVal}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environement) object.Object {

	condition := Eval(node.Confition, env)
	if IsError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if !isTruthy(condition) {
		if node.Alternative != nil {
			return Eval(node.Alternative, env)
		} else {
			return NULL
		}
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environement) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := Builtins[node.Value]; ok {
		return builtin
	}

	return NewError("identifier not fond: " + node.Value)
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

// func evalStatements(stmts []ast.Statement) object.Object {

// 	var obj object.Object

// 	for _, itm := range stmts {
// 		obj = Eval(itm)

// 		if retVal, ok := obj.(*object.ReturnValue); ok {
// 			return retVal.Value
// 		}
// 	}

// 	return obj
// }

func evalProgram(stmts []ast.Statement, env *object.Environement) object.Object {

	var obj object.Object

	for _, itm := range stmts {
		obj = Eval(itm, env)

		// if retVal, ok := obj.(*object.ReturnValue); ok {
		// 	return retVal.Value
		// }
		switch result := obj.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return obj
}

func evalBlockStatements(stmts []ast.Statement, env *object.Environement) object.Object {
	var ret object.Object

	for _, stmt := range stmts {
		ret = Eval(stmt, env)

		if ret != nil {
			retType := ret.Type()
			if retType == object.RETURN_VALUE_OBJ || retType == object.ERROR_OBJ {
				return ret
			}
		}
	}

	return ret
}

func nativeBooltoToBooleanObject(input bool) object.Object {

	if input {
		return TRUE
	}
	return FALSE
}

func NewError(format string, para ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, para...)}
}

func IsError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}

	return false
}
