package main

func f2() { // want "incorrect declaration order"
}

// c - c // want "incorrect declaration order"
const c = "c"

// c2 - c2 // want "incorrect declaration order"
const (
	c2 = "c2"
)
