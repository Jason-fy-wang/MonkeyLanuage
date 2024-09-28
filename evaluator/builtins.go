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
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return NewError("argument to `len` not support. got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got %d, want 1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("argument to `first` must be ARRAY. got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got %d, want 1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length >= 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments. got %d, want 1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElms := make([]object.Object, length-1, length-1)
				copy(newElms, arr.Elements[1:length])
				return &object.Array{Elements: newElms}
			}
			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError("wrong number of arguments. got %d, want 2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return NewError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElems := make([]object.Object, length+1, length+1)
			copy(newElems, arr.Elements)
			newElems[length] = args[1]
			return &object.Array{Elements: newElems}
		},
	},
}
