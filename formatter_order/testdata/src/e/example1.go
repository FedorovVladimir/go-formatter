package main

import "fmt"

var helloText = "hello" // want "formatter_order"

const name = "Vova" // want "formatter_order"

func main() { // want "formatter_order"
	fmt.Println(helloText, name)
}
