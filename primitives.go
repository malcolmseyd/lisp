package main

import "fmt"

func ConsPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)
	if args == Nil {
		panic("cons takes 2 arguments")
	}
	pair, ok := args.(*Pair)
	if !ok {
		panic("bug: cons pair 1")
	}
	left := Car(pair)
	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("bug: cons pair 2")
	}
	if Cdr(pair) != Nil {
		panic(fmt.Sprintf("cons takes 2 arguments, ignored 3rd arg: %#v", Car(pair)))
	}
	right := Car(pair)
	return Cons(left, right)
}

func QuotePrim(o Obj, _ *Env) Obj {
	pair, ok := o.(*Pair)
	if !ok {
		panic("quote takes 1 argument")
	}
	return Car(pair)
}

func EvalPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)
	pair, ok := args.(*Pair)
	if !ok {
		panic("eval takes 1 argument")
	}
	return Eval(Car(pair), e)
}
