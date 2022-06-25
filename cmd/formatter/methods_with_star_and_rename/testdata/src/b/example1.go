package main

type car struct{}

func (c *car) run() {}

func (e car) stop() {} // want "methods_with_star_and_rename"
