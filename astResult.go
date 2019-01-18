package main

import (
	"bytes"
	"fmt"
)

//astResult contains the result of the AST analysis
type astResult struct {
	//array of all variable declarations
	variablesDeclarations []variableDescription
}

type variableDescription struct {
	//type name, in a short version, ex: int, float, Foo...
	typeName string
	//variable name
	varName string
}

//String returns the value of results as a formatted string
func (v *variableDescription) String() string {
	return fmt.Sprintf("{%s}{%s}", v.varName, v.typeName)
}

//String returns the value of results as a formatted string
func (r *astResult) String() string {
	var buffer bytes.Buffer
	for _, e := range r.variablesDeclarations {
		buffer.WriteString(fmt.Sprintf("%s\n", e.String()))
	}

	return buffer.String()
}
