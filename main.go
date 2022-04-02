package main

import (
	"bufio"
	"os"
)

// this function is currently for testing purposes
// this will be the case until we get the parser going
func main() {
	env := MakeEnv(nil)
	BindGlobals(env)

	Print(Eval(Read(bufio.NewReader(os.Stdin)), env))
}
