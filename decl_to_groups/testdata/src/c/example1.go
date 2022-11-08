package main

var ( // want "incorrect single declaration style"
	v1 = "v1"
)
var v2 = "v2"

var v3 = "v3" // want "incorrect single declaration style"
var v4 = "v4"

var v5 = "v5" // want "incorrect single declaration style"
var (
	v6 = "v6"
)
