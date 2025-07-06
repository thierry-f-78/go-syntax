package test

import "fmt"

func shortVarDeclExamples() {
	// Detected cases (forbidden)
	x := 42                  // Short var decl forbidden
	name := "test"           // Short var decl forbidden
	result := someFunction() // Short var decl forbidden

	// Allowed cases
	var y int
	y = 42

	var message string
	message = "test"

	var output string
	output = someFunction()

	// Test with type switch (should be allowed)
	var i interface{}
	i = "hello"

	switch v := i.(type) {
	case string:
		fmt.Println("String:", v)
	case int:
		fmt.Println("Int:", v)
	default:
		fmt.Println("Unknown type")
	}

	// Multiple assignment (forbidden)
	a, b := 1, 2

	// Range with short var (forbidden)
	items := []string{"a", "b", "c"}
	for i, item := range items {
		fmt.Println(i, item)
	}
}

func someFunction() string {
	return "result"
}
