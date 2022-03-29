package main

// this function is currently for testing purposes
// this will be the case until we get the parser going
func main() {
	a := Intern("a")
	b := Intern("b")
	c := Intern("c")
	_, _, _ = a, b, c

	Print(Nil)
	Print(Intern("sym"))
	Print(Cons(a, Nil))
	Print(Cons(a, Cons(b, Nil)))
	Print(Cons(a, Cons(b, Cons(c, Nil))))
	Print(Cons(a, b))
	Print(Cons(a, Cons(b, c)))
}
