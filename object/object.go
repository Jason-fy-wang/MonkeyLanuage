package object

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

func NewEnvironment() *Environement {
	env := &Environement{
		store: make(map[string]Object),
	}
	return env
}

type Environement struct {
	store map[string]Object
}

func (e *Environement) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environement) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
