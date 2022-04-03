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
	if Nil.Equal(o) {
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
		return ApplyProcedure(proc, Evlis(args, e), e)
	default:
		panic(fmt.Sprintf("unknown procedure type %T", proc))
	}
}

func ApplyProcedure(proc *Procedure, argsList Obj, e *Env) Obj {
	args := listToSlice(argsList)
	argsSyms := proc.args
	if len(args) < len(argsSyms) {
		panic(fmt.Sprintf("this procedure takes %v arguments, but was given %v", len(argsSyms), len(args)))
	}

	bodyScope := MakeEnv(proc.scope)
	for i, argSym := range argsSyms {
		bodyScope.Bind(&argSym, args[i])
	}

	if proc.variadic != nil {
		rest := args[len(argsSyms):]
		bodyScope.Bind(proc.variadic, sliceToList(rest))
	} else {
		if len(args) != len(argsSyms) {
			panic(fmt.Sprintf("this procedure takes %v arguments, but was given %v", len(argsSyms), len(args)))
		}
	}

	return Eval(proc.body, bodyScope)
}
