package evaluator

import (
	"testing"

	"com.lanuage/monkey/lexer"
	"com.lanuage/monkey/object"
	"com.lanuage/monkey/parser"
)

func TestEvalIntegetExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, itm := range tests {
		obj := testEval(itm.input)

		testIntegerObject(t, obj, itm.expected)
	}

}

func TestEvalBoolExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, itm := range tests {
		obj := testEval(itm.input)
		testBoolObject(t, obj, itm.expect)
	}
}

func TestBangOperation(t *testing.T) {
	tests := []struct {
		input  string
		expect bool
	}{
		{"!5", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, itm := range tests {
		obj := testEval(itm.input)
		testBoolObject(t, obj, itm.expect)
	}
}

func TestEvalPreMinusOperation(t *testing.T) {
	tests := []struct {
		intput string
		expect int64
	}{
		{"5", 5},
		{"-5", -5},
		{"10", 10},
		{"-10", -10},
	}

	for _, itm := range tests {
		obj := testEval(itm.intput)
		testIntegerObject(t, obj, itm.expect)
	}
}

func testEval(input string) object.Object {

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParserProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expect int64) bool {

	results, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("expect object.Integer, got: %v", obj)
		return false
	}

	if results.Value != expect {
		t.Errorf("expect %d, got %v", expect, results.Value)
		return false
	}

	return true
}

func testBoolObject(t *testing.T, obj object.Object, expect bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("expect bool, got %v", result)
		return false
	}

	if result.Value != expect {

		t.Errorf("expect value: %t, got %t", expect, result.Value)
		return false
	}
	return true
}
