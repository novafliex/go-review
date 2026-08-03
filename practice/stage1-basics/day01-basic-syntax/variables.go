package main

import "fmt"

// variablesFunc demonstrates common ways to declare variables and constants in Go.
func variablesFunc() {
	// Go provides three common ways to declare variables.
	var a int = 1
	var b = 2
	c := 4

	fmt.Println(a, b, c)

	// Declare multiple variables in a single declaration.
	var a1, a2, a3 = 1, 2, 3

	// A short variable declaration can contain values of different types.
	b1, b2, b3 := 1, "a", 1.2

	fmt.Println(a1, a2, a3)
	fmt.Println(b1, b2, b3)

	// The := operator(操作) can only be used inside functions(函数内部).
	// At least one variable on the left-hand side must be new.
	name := "Mike"
	name, age := "Bob", 18

	fmt.Println(name, age)

	// A constant expression cannot depend on variables declared with var or :=.
	// Constants can depend on other constants and constant expressions.
	const d = 2
	const d1, d2 = 1, 2
	const d3 = d1 + d2 + 1

	fmt.Println(d)
	fmt.Println(d3)

	// Grouped constant declaration.
	const (
		e  = 10
		e1 = "a"
		e2 = 1.23
	)

	fmt.Println(e, e1, e2)

	// iota starts at 0, increments for each constant specification(规格),
	// and resets to 0 in each new const block.
	// It increments according to the declaration position,
	// even when the current(当前的) line does not explicitly(明确地) use iota.
	const (
		f  = 0    // iota = 0
		f1 = 2    // iota = 1
		f2 = iota // iota = 2
		f3        // repeats the previous expression: iota = 3
		_         // discards iota = 4
		f4        // repeats the previous expression: iota = 5
	)

	fmt.Println(f, f1, f2, f3, f4)
}
