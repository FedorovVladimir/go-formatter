package main

import ( // want "rm decl" "incorrect single declaration style"
	"fmt"
)

var ( // want "rm decl" "incorrect single declaration style"
	v1 = "v1"
)

const ( // want "rm decl" "incorrect single declaration style"
	c1 = "c1"
)

func main() {
	fmt.Println()
}
