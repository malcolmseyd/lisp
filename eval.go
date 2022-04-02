package main

import "fmt"

func Eval(o Obj, e *Env) Obj {
	switch o := o.(type) {
	case *Primitive, *Num:
		return o
	case *Symbol:
		return e.Resolve(o)
	case *Pair:
		return Apply(Eval(Car(o), e), Cdr(o), e)
	default:
		panic(fmt.Sprintf("unknown object %#v passed to eval", o))
	}
}

func Evlis(o Obj, e *Env) Obj {
	if o == Nil {
		return Nil
	}
	if pair, ok := o.(*Pair); ok {
		return Cons(Eval(Car(pair), e), Evlis(Cdr(pair), e))
	}
	panic("evlist called on a non-list object %#v")
}

func Apply(proc Obj, args Obj, e *Env) Obj {
	switch proc := proc.(type) {
	case Primitive:
		return proc(args, e)
	default:
		panic(fmt.Sprintf("unknown procedure type %T", proc))
	}
}
