package main

var ( // want "rm decl"
	v1 = "v1"
)
var v2 = "v2" // want "rm decl" "incorrect single declaration style"
