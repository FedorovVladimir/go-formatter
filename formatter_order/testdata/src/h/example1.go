package main

import "fmt"

func main() { // want "incorrect declaration order"
	var a int
	fmt.Println(a)
}

func Hello() { // want "incorrect declaration order"

}

func NewHello1() { // want "incorrect declaration order"

}

func newHello2() { // want "incorrect declaration order"

}
