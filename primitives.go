package main

import (
	"fmt"
	"os"
)

func LambdaPrim(o Obj, e *Env) Obj {
	pair, ok := o.(*Pair)
	if !ok {
		panic("lambda takes 2 arguments")
	}
	args := Car(pair)

	argsSyms := []Symbol{}
	for !Nil.Equal(args) {
		argList, ok := args.(*Pair)
		if !ok {
			panic("lambda args must be an argument list")
		}
		sym, ok := Car(argList).(*Symbol)
		if !ok {
			panic("arguments must be symbols")
		}
		argsSyms = append(argsSyms, *sym)

		args = Cdr(argList)
	}

	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("lambda takes 2 arguments")
	}
	body := Car(pair)

	return MakeProcedure(argsSyms, body, e)
}

func EqPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)

	pair, ok := args.(*Pair)
	if !ok {
		panic("eq takes 2 arguments")
	}
	v1 := Car(pair)

	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("eq takes 2 arguments")
	}
	v2 := Car(pair)

	if v1.Type() != v2.Type() {
		return Nil
	}

	switch v1 := v1.(type) {
	case *Symbol:
		return boolToLisp(*v1 == *v2.(*Symbol))
	case *Num:
		return boolToLisp(*v1 == *v2.(*Num))
	default:
		// reference equality for misc
		return boolToLisp(v1 == v2)
	}
}

func LessPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)

	pair, ok := args.(*Pair)
	if !ok {
		panic("eq takes 2 arguments")
	}
	v1, ok := Car(pair).(*Num)
	if !ok {
		panic("args should be numbers")
	}

	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("eq takes 2 arguments")
	}
	v2, ok := Car(pair).(*Num)
	if !ok {
		panic("args should be numbers")
	}

	return boolToLisp(v1.n < v2.n)
}

func ConsPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)
	if Nil.Equal(args) {
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
	if !Nil.Equal(Cdr(pair)) {
		panic(fmt.Sprintf("cons takes 2 arguments, ignored 3rd arg: %#v", Car(pair)))
	}
	right := Car(pair)
	return Cons(left, right)
}

func CarPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)
	if Nil.Equal(args) {
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
	if Nil.Equal(args) {
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
	if Nil.Equal(o) {
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
	if Nil.Equal(o) {
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
	if Nil.Equal(o) {
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

	if !Nil.Equal(test) {
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

func ApplyPrim(o Obj, e *Env) Obj {
	args := Evlis(o, e)
	pair, ok := args.(*Pair)
	if !ok {
		panic("apply takes 2 arguments")
	}
	proc := Car(pair)
	pair, ok = Cdr(pair).(*Pair)
	if !ok {
		panic("apply takes 2 arguments")
	}
	return Apply(proc, Car(pair), e)

}

func AddPrim(o Obj, e *Env) Obj {
	acc := int64(0)
	args := Evlis(o, e)
	pair, ok := args.(*Pair)
	if !ok {
		if Nil.Equal(args) {
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
		if Nil.Equal(Cdr(pair)) {
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
		if Nil.Equal(args) {
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
		if Nil.Equal(Cdr(pair)) {
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
		if Nil.Equal(args) {
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
		if Nil.Equal(Cdr(pair)) {
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
		if Nil.Equal(args) {
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
		if Nil.Equal(Cdr(pair)) {
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
		if Nil.Equal(args) {
			os.Exit(0)
		}
		panic("eval takes 1 argument")
	}
	if n, ok := Car(pair).(*Num); ok {
		os.Exit(int(n.n))
	}
	panic("exit takes an number")
}
