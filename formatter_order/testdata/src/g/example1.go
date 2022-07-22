package main

import "fmt"

func main() { // want "incorrect declaration order"
	var a int
	fmt.Println(a)
}

func Hello() { // want "incorrect declaration order"

}
