package main

import "fmt"

var helloText = "hello" // want "formatter_order"

const name = "Vova" // want "formatter_order"

func main() {
	fmt.Println(helloText, name)
}
