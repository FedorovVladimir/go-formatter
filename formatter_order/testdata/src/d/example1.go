package main

var v1 = "v1"

type e struct{} // want "incorrect declaration order"

var v2 = "v2" // want "incorrect declaration order"
