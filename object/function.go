package object

import (
	"bytes"
	"strings"

	"com.language/monkey/ast"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatements
	Env        *Environement
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}

	for _, itm := range f.Parameters {
		params = append(params, itm.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
