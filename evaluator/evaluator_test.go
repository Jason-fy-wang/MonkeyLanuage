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

		{"1<2", true},
		{"1>2", false},
		{"1<1", false},
		{"1>1", false},
		{"1==1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},

		{"true==true", true},
		{"false==false", true},
		{"false == true", false},
		{"true!=false", true},
		{"false!=true", true},
		{"(1<2) == true", true},
		{"(1<2) == false", false},
		{"(1>2) == true", false},
		{"(1>2) == false", true},
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

		{"5+5+5+5-10", 10},
		{"2*2*2*2*2", 32},
		{"-50+100+ -50", 0},
		{"5*2+10", 20},
		{"20+2*10", 40},
		{"20+ 2* -10", 0},
		{"50/2*2+10", 60},
		{"2*(5+10)", 30},
		{"3*3*3+10", 37},
		{"3*(3*3)+10", 37},
		{"(5+10*2+15/3)*2+ -10", 50},
	}

	for _, itm := range tests {
		obj := testEval(itm.intput)
		testIntegerObject(t, obj, itm.expect)
	}
}

func TestIfParseExpression(t *testing.T) {
	tests := []struct {
		input  string
		expect interface{}
	}{
		{"if(false) {10}", nil},
		{"if(true) {10}", 10},
		{"if(1) {10}", 10},
		{"if(1<2) {10}", 10},
		{"if(1>2) {10}", nil},
		{"if(1>2) {10} else {20}", 20},
		{"if(1<2) {10} else{20}", 10},
	}

	for _, itm := range tests {
		obj := testEval(itm.input)
		integer, ok := itm.expect.(int)
		if !ok {
			testNullObject(t, obj)
		} else {
			testIntegerObject(t, obj, int64(integer))
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got %T (%+v)", obj, obj)
		return false
	}
	return true
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
