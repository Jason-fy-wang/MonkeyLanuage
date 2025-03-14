package evaluator

import (
	"testing"

	"com.language/monkey/lexer"
	"com.language/monkey/object"
	"com.language/monkey/parser"
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

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknow operator:-BOOLEAN",
		},
		{
			"true + false",
			"unknow operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknow operator: BOOLEAN + BOOLEAN",
		},
		{
			"if(10>1){true +  false;}",
			"unknow operator: BOOLEAN + BOOLEAN",
		},
		{
			"if(10>1){ if (10 > 1) { return true +  false;}  return 1;}",
			"unknow operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not fond: foobar",
		},
		{
			`"hello" - "world"`,
			"unknow operator:STRING - STRING",
		},
	}

	for _, itm := range tests {
		obj := testEval(itm.input)

		err, ok := obj.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got %T(%+v)", obj, obj)
			continue
		}

		if err.Message != itm.expected {
			t.Errorf("expect error mesage: %s, got %s", itm.expected, err.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {

	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a+b; c;", 10},
	}

	for _, itm := range tests {
		testIntegerObject(t, testEval(itm.input), itm.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) {x+2;};"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not function. got %T (%-v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("wrong parameters count. Parameters=%-v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("expect parameter x, got %-v", fn.Parameters[0].String())
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("expect body: %s, got %s", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		intput   string
		expected int64
	}{
		{
			"let add = fn(x,y){x+y;};  add(5+5, add(5, 5));",
			20,
		},
		{
			"let identity=fn(x) {x;};  identity(5);",
			5,
		},
		{
			"let identity=fn(x) {return x;};  identity(5);",
			5,
		},
		{
			"let double=fn(x){x*2;};  double(5);",
			10,
		},
		{
			"let add = fn(x,y){x+y;};   add(5,5);",
			10,
		},

		{
			"fn(x){x;}(5)",
			5,
		},
	}

	for _, itm := range tests {
		evaluated := testEval(itm.intput)
		testIntegerObject(t, evaluated, itm.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world"`

	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("expect object.string. got %v", evaluated)
	}

	if str.Value != "hello world" {
		t.Fatalf("expect %s, got %s", input, str)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"hello" + " " + "world"`

	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not string. got %T", evaluated)
	}

	if str.Value != "hello world" {
		t.Errorf("String has wrong value. got %s", str)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not support. got INTEGER"},
		{`len("one","two")`, "wrong number of parameters. got 2, want 1"},
	}

	for _, itm := range tests {
		evaluated := testEval(itm.input)

		switch expected := itm.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1,2*2, 3+3]"

	evaluated := testEval(input)

	arrObj, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object not array. got %T (%+v)", evaluated, evaluated)
	}

	if len(arrObj.Elements) != 3 {
		t.Fatalf("expect elements number is 3,got %d", len(arrObj.Elements))
	}

	testIntegerObject(t, arrObj.Elements[0], 1)
	testIntegerObject(t, arrObj.Elements[1], 4)
	testIntegerObject(t, arrObj.Elements[2], 6)

}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1,2,3][0]",
			1,
		},
		{
			"[1,2,3][2]",
			3,
		},
		{
			"[1,2,3][1]",
			2,
		},
		{
			"[1,2,3][-1]",
			nil,
		},
		{
			"[1,2,3][4]",
			nil,
		},
	}

	for _, itm := range tests {
		evaluated := testEval(itm.input)
		integer, ok := itm.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}

	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
		{
			"one":10-9,
			two: 1+1,
			"thr"+"ee" : 6/2,
			4:4,
			true:5,
			false: 6
		}
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't result HASH.got %T (%v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. got %d", len(result.Pairs))
	}
	for keyExpect, valExpect := range expected {
		pair, ok := result.Pairs[keyExpect]
		if !ok {
			t.Errorf("no pair for given key : %v", keyExpect)
		}

		testIntegerObject(t, pair.Value, valExpect)
	}
}

func TestHashIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"];`,
			5,
		},
		{
			`{"foo": 5}["BAR"];`,
			NULL,
		},
		{
			`{}["BAR"];`,
			NULL,
		},
		{
			`{true: 5}[true];`,
			5,
		},
		{
			`{false: 5}[false];`,
			5,
		},
		{
			`let key ="foo"; {"foo": 5}[key];`,
			5,
		},
	}
	for _, itm := range tests {
		evaluated := testEval(itm.input)
		expected, ok := itm.expected.(int)
		if !ok {
			testNullObject(t, evaluated)
		} else {
			testIntegerObject(t, evaluated, int64(expected))
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
	env := object.NewEnvironment()
	return Eval(program, env)
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
