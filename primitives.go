package main

import (
	"fmt"
	"os"
)

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
		panic("cons takes 2 arguments")
	}
	if Cdr(pair) != Nil {
		panic(fmt.Sprintf("cons takes 2 arguments, ignored 3rd arg: %#v", Car(pair)))
	}
	right := Car(pair)
	return Cons(left, right)
}

func CarPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)
	if args == Nil {
		panic("car takes 1 argument")
	}

	pair, ok := args.(*Pair)
	if !ok {
		panic("bug: car pair")
	}
	pair, ok = Car(pair).(*Pair)
	if !ok {
		panic("car takes pairs as arguments")
	}
	return Car(pair)
}

func CdrPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)
	if args == Nil {
		panic("car takes 1 argument")
	}

	pair, ok := args.(*Pair)
	if !ok {
		panic("bug: car pair")
	}
	pair, ok = Car(pair).(*Pair)
	if !ok {
		panic("car takes pairs as arguments")
	}
	return Cdr(pair)
}

func DefinePrim(o Obj, e *Env) Obj {
	if o == Nil {
		panic("define takes 2 arguments")
	}

	pair, ok := o.(*Pair)
	if !ok {
		panic("define bug 1")
	}
	name, ok := Car(pair).(*Symbol)
	if !ok {
		panic("the first argument to define is a symbol")
	}
	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("define takes 2 arguments")
	}

	e.Bind(name, Eval(Car(pair), e))
	return Nil
}

func SetPrim(o Obj, e *Env) Obj {
	if o == Nil {
		panic("set takes 2 arguments")
	}

	pair, ok := o.(*Pair)
	if !ok {
		panic("set bug 1")
	}
	name, ok := Car(pair).(*Symbol)
	if !ok {
		panic("the first argument to set is a symbol")
	}
	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("set takes 2 arguments")
	}

	e.Set(name, Eval(Car(pair), e))
	return Nil
}

func IfPrim(o Obj, e *Env) Obj {
	if o == Nil {
		panic("if takes 3 arguments")
	}
	pair, ok := o.(*Pair)
	if !ok {
		panic("if bug 1")
	}
	test := Eval(Car(pair), e)

	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("if takes 3 arguments")
	}
	expr1 := Car(pair)

	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("if takes 3 arguments")
	}
	expr2 := Car(pair)

	if test != Nil {
		return Eval(expr1, e)
	}
	return Eval(expr2, e)
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

func AddPrim(o Obj, e *Env) Obj {
	acc := int64(0)
	args := Evlis(o, e)
	pair, ok := args.(*Pair)
	if !ok {
		if args == Nil {
			return MakeNum(acc)
		}
		panic("bug in addprim")
	}
	for {
		n, ok := Car(pair).(*Num)
		if !ok {
			panic("+ only takes number arguments")
		}
		acc += n.n
		if Cdr(pair) == Nil {
			break
		}
		pair, ok = Cdr(pair).(*Pair)
		if !ok {
			panic("bad args to + special form")
		}
	}
	return MakeNum(acc)
}

func SubPrim(o Obj, e *Env) Obj {
	acc := int64(0)
	first := true
	args := Evlis(o, e)
	pair, ok := args.(*Pair)
	if !ok {
		if args == Nil {
			return MakeNum(acc)
		}
		panic("bug in subprim")
	}
	for {
		n, ok := Car(pair).(*Num)
		if !ok {
			panic("- only takes number arguments")
		}
		if first {
			acc = n.n
		} else {
			acc -= n.n
		}
		if Cdr(pair) == Nil {
			break
		}
		pair, ok = Cdr(pair).(*Pair)
		if !ok {
			panic("bad args to - special form")
		}
		first = false
	}
	// unary minus is negation
	if first {
		return MakeNum(-acc)
	}
	return MakeNum(acc)
}

func MulPrim(o Obj, e *Env) Obj {
	acc := int64(1)
	args := Evlis(o, e)
	pair, ok := args.(*Pair)
	if !ok {
		if args == Nil {
			return MakeNum(acc)
		}
		panic("bug in multprim")
	}
	for {
		n, ok := Car(pair).(*Num)
		if !ok {
			panic("/ only takes number arguments")
		}
		acc *= n.n
		if Cdr(pair) == Nil {
			break
		}
		pair, ok = Cdr(pair).(*Pair)
		if !ok {
			panic("bad args to / special form")
		}
	}
	return MakeNum(acc)
}

func DivPrim(o Obj, e *Env) Obj {
	acc := int64(1)
	first := true
	args := Evlis(o, e)
	pair, ok := args.(*Pair)
	if !ok {
		if args == Nil {
			return MakeNum(acc)
		}
		panic("bug in multprim")
	}
	for {
		n, ok := Car(pair).(*Num)
		if !ok {
			panic("/ only takes number arguments")
		}
		if first {
			acc = n.n
		} else {
			acc /= n.n
		}
		if Cdr(pair) == Nil {
			break
		}
		pair, ok = Cdr(pair).(*Pair)
		if !ok {
			panic("bad args to / special form")
		}
		first = false
	}
	return MakeNum(acc)
}

func ExitPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)
	pair, ok := args.(*Pair)
	if !ok {
		if args == Nil {
			os.Exit(0)
		}
		panic("eval takes 1 argument")
	}
	if n, ok := Car(pair).(*Num); ok {
		os.Exit(int(n.n))
	}
	panic("exit takes an number")
}
