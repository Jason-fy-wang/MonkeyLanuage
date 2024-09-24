package evaluator

import "com.language/monkey/object"

var Builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of parameters. got %d, want 1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			default:
				return NewError("argument to `len` not support. got %s", args[0].Type())
			}
		},
	},
}
