package object

type ReturnValue struct {
	Value Object
}

func (rb *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rb *ReturnValue) Inspect() string {
	return rb.Value.Inspect()
}
