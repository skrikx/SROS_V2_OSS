package main

import (
	"os"
)

func main() {
	os.Exit(Execute(os.Args[1:], RunOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}))
}
