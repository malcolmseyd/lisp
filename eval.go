package main

import "fmt"

func Eval(o Obj, e *Env) Obj {
	switch o := o.(type) {
	case *Primitive, *Procedure, *Num:
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
	case *Procedure:
		args := Evlis(args, e)

		bodyScope := MakeEnv(proc.scope)
		argsSyms := proc.args

		var argsList *Pair = nil
		var ok bool
		for i, argSym := range argsSyms {
			if args == Nil {
				panic(fmt.Sprintf("this procedure takes %v arguments, but was only given %v", len(argsSyms), i))
			}
			argsList, ok = args.(*Pair)
			if !ok {
				panic("bug: args should be a list")
			}
			bodyScope.Bind(&argSym, Car(argsList))
			args = Cdr(argsList)
		}

		return Eval(proc.body, bodyScope)
	default:
		panic(fmt.Sprintf("unknown procedure type %T", proc))
	}
}
