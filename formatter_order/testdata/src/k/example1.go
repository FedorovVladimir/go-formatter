package main

func main() { // want "incorrect declaration order"
}

// v - v // want "incorrect declaration order"
var v = "v"

// v2 - v2 // want "incorrect declaration order"
var (
	v2 = "v2"
)
