package main

import (
	"os"
	"queue/internals"
)

func main() {
	arg := os.Args
	c := internals.MakeConsumer(arg[0])
	c.Subscribe()
}
