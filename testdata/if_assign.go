package test

import (
	"errors"
	"fmt"
)

func ifAssignExamples() {
	// Detected cases (forbidden)
	if err = someErrorFunction(); err != nil {
		fmt.Println("Error:", err)
	}

	if result = computeResult(); result != nil {
		fmt.Println("Result:", result)
	}

	if value = getValue(); value != nil {
		fmt.Println("Value:", value)
	}

	// Allowed cases
	var err error
	err = someErrorFunction()
	if err != nil {
		fmt.Println("Error:", err)
	}

	var result interface{}
	result = computeResult()
	if result != nil {
		fmt.Println("Result:", result)
	}

	// If with short var (other rule, but also forbidden)
	if x := getValue(); x != nil {
		fmt.Println("X:", x)
	}

	// If with assignment but no nil check (allowed)
	if count = getCount(); count > 0 {
		fmt.Println("Count:", count)
	}

	// If without assignment (allowed)
	if someCondition() {
		fmt.Println("Condition met")
	}
}

var err error
var result interface{}
var value interface{}
var count int

func someErrorFunction() error {
	return errors.New("test error")
}

func computeResult() interface{} {
	return "result"
}

func getValue() interface{} {
	return "value"
}

func getCount() int {
	return 5
}

func someCondition() bool {
	return true
}
