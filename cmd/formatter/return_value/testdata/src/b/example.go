package a

func returnValue(a, b int) (e int, d bool) { // want "return value"
	return a + b, a > b
}
