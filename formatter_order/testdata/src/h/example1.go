package main

import "fmt"

func main() { // want "formatter_order"
	var a int
	fmt.Println(a)
}

func Hello() {

}

func NewHello1() { // want "formatter_order"

}

func newHello2() { // want "formatter_order"

}
