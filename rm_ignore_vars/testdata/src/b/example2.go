package main

var _, b = "b", "b"

var _, _ = "b", "b" // want "rm_ignore_vars"
