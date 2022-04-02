package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// this function is currently for testing purposes
// this will be the case until we get the parser going
func main() {
	e := MakeEnv(nil)
	BindGlobals(e)

	r := bufio.NewReader(os.Stdin)

	for {
		repl(r, e)
	}
}

func repl(r io.RuneScanner, e *Env) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic caught:", r)
		}
	}()

	fmt.Print("> ")
	Print(Eval(Read(r), e))
}
