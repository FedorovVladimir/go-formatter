package main

// doc for file

import "fmt"

func main() { // want "incorrect declaration order"
	fmt.Println("Hello")
}

func Hello() { // want "incorrect declaration order"

}
