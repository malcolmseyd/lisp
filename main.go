package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"os"
)

var _, noREPL = os.LookupEnv("NO_REPL")

// this function is currently for testing purposes
// this will be the case until we get the parser going
func main() {
	e := MakeEnv(nil)
	BindGlobals(e)

	loadPrelude(e)

	r := bufio.NewReader(os.Stdin)

	for {
		repl(r, e)
	}
}

//go:embed prelude.lisp
var prelude string

func loadPrelude(e *Env) {
	defer func() {
		if r := recover(); r != nil && r != "EOF" {
			fmt.Println("panic caught loading stdlib:", r)
		}
	}()
	s := bufio.NewReader(bytes.NewBufferString(prelude))

	// go until we panic (EOF panics)
	for {
		Eval(Read(s), e)
	}
}

func repl(r *bufio.Reader, e *Env) {
	defer func() {
		if r := recover(); r != nil {
			if r == "EOF" {
				os.Exit(0)
			}
			fmt.Println("panic caught:", r)
		}
	}()

	if noREPL {
		Print(Eval(Read(r), e))
	} else {
		fmt.Print("> ")
		result := Eval(Read(r), e)
		fmt.Print("< ")
		Print(result)
	}
}
