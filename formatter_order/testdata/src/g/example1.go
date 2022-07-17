package main

import "fmt"

func main() { // want "formatter_order"
	var a int
	fmt.Println(a)
}

func Hello() { // want "formatter_order"

}
