package object

import (
	"bytes"
	"strings"
)

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}

	for _, itm := range ao.Elements {
		elements = append(elements, itm.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
