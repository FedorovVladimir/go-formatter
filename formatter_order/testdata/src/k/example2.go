package main

func f() { // want "formatter_order"
}

// s - s // want "formatter_order"
type s struct {
}

// s2 - s2 // want "formatter_order"
type (
	s2 struct {
	}
)
