package main

type Operation int

const (
	Add Operation = iota // want "start_enums_at_one"
	Subtract
	Multiply
)
