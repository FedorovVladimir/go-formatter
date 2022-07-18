package main

// doc for file

import "fmt"

func main() { // want "formatter_order"
	fmt.Println("Hello")
}

func Hello() { // want "formatter_order"

}
