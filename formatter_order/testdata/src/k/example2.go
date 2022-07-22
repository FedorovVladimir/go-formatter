package main

func f() { // want "incorrect declaration order"
}

// s - s // want "incorrect declaration order"
type s struct {
}

// s2 - s2 // want "incorrect declaration order"
type (
	s2 struct {
	}
)
