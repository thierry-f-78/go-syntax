package test

import (
	"fmt"
	"os"
)

func mixedExamples() {
	// Multiple violations in the same function

	// Short var decl (forbidden)
	filename := "test.txt"

	// If assign with nil check (forbidden)
	if file, err = os.Open(filename); err != nil {
		fmt.Println("Cannot open file:", err)
		return
	}
	defer file.Close()

	// Short var decl in range (forbidden)
	data := []byte("hello world")
	for i, b := range data {
		fmt.Printf("Byte %d: %c\n", i, b)
	}

	// Short var decl in condition (forbidden)
	if content := readFile(filename); content != "" {
		fmt.Println("Content:", content)
	}

	// Nested short var decl (forbidden)
	if true {
		inner := "nested"
		fmt.Println(inner)
	}

	// Function with short var return (forbidden)
	result := processData(data)
	fmt.Println("Result:", result)

	// Type assertion with short var (forbidden)
	var x interface{}
	x = "test"
	if str, ok := x.(string); ok {
		fmt.Println("String:", str)
	}
}

var file *os.File
var err error

func readFile(filename string) string {
	return "file content"
}

func processData(data []byte) string {
	return "processed"
}
