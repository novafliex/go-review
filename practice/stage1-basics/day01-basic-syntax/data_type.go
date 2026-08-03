package main

import "fmt"

// practiceDataTypes demonstrates the common basic data types in Go.
func practiceDataTypes() {
	// Basic data types.
	var integerValue int = 100
	var floatValue float64 = 3.14
	var boolValue bool = false
	var stringValue string = "hello"
	var byteValue byte = 'A'
	var runeValue rune = '你'

	fmt.Println("Basic values:")
	fmt.Println(integerValue, floatValue, boolValue, stringValue)

	// byte is an alias(别名) for uint8.
	// It is commonly used to represent raw(原始) bytes or ASCII characters.
	fmt.Printf("byte: value=%d, character=%c, type=%T\n",
		byteValue, byteValue, byteValue)

	// rune is an alias for int32.
	// It is used to represent a Unicode code point.
	fmt.Printf("rune: value=%d, character=%c, type=%T\n",
		runeValue, runeValue, runeValue)

	// A string is immutable(不可变的) in Go.
	// Convert(转换) it to []byte before modifying(修改) an ASCII character.
	text := "hello"
	bytes := []byte(text)
	bytes[0] = 'H'
	text = string(bytes)

	fmt.Println("Modified ASCII string:", text)

	// Use []rune when modifying Unicode characters.
	chineseText := "你好"
	runes := []rune(chineseText)
	runes[0] = '我'
	chineseText = string(runes)

	fmt.Println("Modified Unicode string:", chineseText)

	// Zero values of basic data types.
	var zeroInt int
	var zeroFloat float64
	var zeroBool bool
	var zeroString string
	var zeroByte byte
	var zeroRune rune

	fmt.Println("Zero values:")
	fmt.Printf(
		"int=%d, float64=%f, bool=%t, string=%q, byte=%d, rune=%d\n",
		zeroInt,
		zeroFloat,
		zeroBool,
		zeroString,
		zeroByte,
		zeroRune,
	)

	// Go requires explicit(明确的) conversion(转换) between different numeric(数字的) types.
	number := 10
	decimal := 2.5
	result := float64(number) + decimal

	fmt.Println("Converted result:", result)

	// Type inference(推断).
	inferredInt := 20
	inferredFloat := 6.28
	inferredString := "Go"
	inferredBool := true

	fmt.Printf(
		"Inferred types: %T, %T, %T, %T\n",
		inferredInt,
		inferredFloat,
		inferredString,
		inferredBool,
	)
}
