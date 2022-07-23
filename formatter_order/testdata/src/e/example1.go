package main

import "fmt"

var helloText = "hello" // want "incorrect declaration order"

const name = "Vova" // want "incorrect declaration order"

func main() {
	fmt.Println(helloText, name)
}
