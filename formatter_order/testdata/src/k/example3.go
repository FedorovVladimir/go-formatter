package main

func f2() { // want "formatter_order"
}

// c - c // want "formatter_order"
const c = "c"

// c2 - c2 // want "formatter_order"
const (
	c2 = "c2"
)
