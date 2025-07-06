package test

import "fmt"

func nolintExamples() {
	// This short var decl will be ignored thanks to rule name
	x := 42 // nolint:short-var-decl

	// This if init will be ignored thanks to rule name
	if err = someFunction(); err != nil { // nolint:if-init
		fmt.Println("Error:", err)
	}

	// Ignore all rules for this line
	y := "test" // nolint:all

	// This short var decl will be ignored thanks to error code
	a := 123 // nolint:GS001

	// This if init will be ignored thanks to error code
	if value = getValue(); value != nil { // nolint:GS002
		fmt.Println("Value:", value)
	}

	// This case will be detected because no nolint
	z := 100

	// This case will be detected because no nolint
	if result = anotherFunction(); result != nil {
		fmt.Println("Result:", result)
	}
}

func getValue() interface{} {
	return "test"
}

var err error
var result interface{}

func someFunction() error {
	return nil
}

func anotherFunction() interface{} {
	return "result"
}
