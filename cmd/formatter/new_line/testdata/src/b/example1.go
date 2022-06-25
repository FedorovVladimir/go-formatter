package main

type car struct{}

func (c *car) run()  {} // want "new_line"
func (c *car) stop() {} // want "new_line"
