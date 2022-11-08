package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		i := i // want 'shadowed_var is forbidden'
		fmt.Println(i)
	}
}
