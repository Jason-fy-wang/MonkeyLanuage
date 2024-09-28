package object

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

func NewEnclosedEnvironment(outer *Environement) *Environement {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environement {
	env := &Environement{
		store: make(map[string]Object),
		outer: nil,
	}
	return env
}

type Environement struct {
	store map[string]Object
	outer *Environement
}

func (e *Environement) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environement) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
