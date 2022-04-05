package main

import (
	"os"
)

func LambdaPrim(o Obj, e *Env) Obj {
	formArgs := listToSlice(o)
	if len(formArgs) < 2 {
		panic("lambda takes at least 2 arguments")
	}
	args, variadic := improperListToSlice(formArgs[0])

	variadicSym, ok := variadic.(*Symbol)
	if variadic != nil && !ok {
		panic("arguments must be symbols")
	}

	argsSyms := []Symbol{}
	for _, arg := range args {
		sym, ok := arg.(*Symbol)
		if !ok {
			panic("arguments must be symbols")
		}
		argsSyms = append(argsSyms, *sym)
	}

	body := sliceToList(formArgs[1:])

	if variadicSym != nil {
		return MakeVariadicProcedure(argsSyms, *variadicSym, body, e)
	}
	return MakeProcedure(argsSyms, body, e)
}

func EqPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))

	v1 := args[0]
	v2 := args[1]

	if v1.Type() != v2.Type() {
		return Nil
	}

	switch v1 := v1.(type) {
	case *Symbol:
		return boolToLisp(v1.Equal(v2))
	case *Num:
		return boolToLisp(*v1 == *v2.(*Num))
	default:
		// reference equality for misc
		return boolToLisp(v1 == v2)
	}
}

func LessPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))

	v1, ok := args[0].(*Num)
	if !ok {
		panic("args should be numbers")
	}

	v2, ok := args[1].(*Num)
	if !ok {
		panic("args should be numbers")
	}

	return boolToLisp(v1.n < v2.n)
}

func ConsPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 2 {

		panic("cons takes 2 arguments")
	}
	left := args[0]
	right := args[1]
	return Cons(left, right)
}

func CarPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("car takes 1 argument")
	}

	pair, ok := args[0].(*Pair)
	if !ok {
		panic("car takes pairs as arguments")
	}
	return Car(pair)
}

func CdrPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("cdr takes 1 argument")
	}

	pair, ok := args[0].(*Pair)
	if !ok {
		panic("cdr takes pairs as arguments")
	}
	return Cdr(pair)
}

func DefinePrim(o Obj, e *Env) Obj {
	args := listToSlice(o)
	if len(args) != 2 {
		panic("define takes 2 arguments")
	}

	name, ok := args[0].(*Symbol)
	if !ok {
		panic("the first argument to define is a symbol")
	}

	expr := args[1]
	return e.Bind(name, Eval(expr, e))
}

func SetPrim(o Obj, e *Env) Obj {
	args := listToSlice(o)
	if len(args) != 2 {
		panic("set takes 2 arguments")
	}

	name, ok := args[0].(*Symbol)
	if !ok {
		panic("the first argument to set is a symbol")
	}

	expr := args[1]
	return e.Set(name, Eval(expr, e))
}

func SetCarPrim(o Obj, e *Env) Obj {
	args := listToSlice(o)
	if len(args) != 2 {
		panic("set takes 2 arguments")
	}

	name, ok := args[0].(*Symbol)
	if !ok {
		panic("the first argument to set is a symbol")
	}

	expr := args[1]
	return e.SetCar(name, Eval(expr, e))
}

func SetCdrPrim(o Obj, e *Env) Obj {
	args := listToSlice(o)
	if len(args) != 2 {
		panic("set takes 2 arguments")
	}

	name, ok := args[0].(*Symbol)
	if !ok {
		panic("the first argument to set is a symbol")
	}

	expr := args[1]
	return e.SetCdr(name, Eval(expr, e))
}

func IfPrim(o Obj, e *Env) Obj {
	args := listToSlice(o)
	if len(args) != 3 {
		panic("if takes 3 arguments")
	}
	test := Eval(args[0], e)
	expr1 := args[1]
	expr2 := args[2]
	if !Nil.Equal(test) {
		return Eval(expr1, e)
	}
	return Eval(expr2, e)
}

func QuotePrim(o Obj, _ *Env) Obj {
	args := listToSlice(o)
	if len(args) != 1 {
		panic("quote takes 1 argument")
	}
	return args[0]
}

func EvalPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("eval takes 1 argument")
	}
	return Eval(args[0], e)
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
	args := listToSlice(Evlis(o, e))
	if len(args) == 0 {
		return MakeNum(acc)
	}
	for _, n := range args {
		n, ok := n.(*Num)
		if !ok {
			panic("+ only takes number arguments")
		}

		acc += n.n
	}
	return MakeNum(acc)
}

func SubPrim(o Obj, e *Env) Obj {
	acc := int64(0)
	args := listToSlice(Evlis(o, e))
	if len(args) == 0 {
		return MakeNum(acc)
	}
	for i, n := range args {
		n, ok := n.(*Num)
		if !ok {
			panic("- only takes number arguments")
		}

		// first element is minuend, following are subtrahend (i googled this lol)
		if i == 0 {
			acc = n.n
		} else {
			acc -= n.n
		}
	}
	// special case: unary minus is negation
	if len(args) == 1 {
		return MakeNum(-acc)
	}
	return MakeNum(acc)
}

func MulPrim(o Obj, e *Env) Obj {
	acc := int64(1)
	args := listToSlice(Evlis(o, e))
	if len(args) == 0 {
		return MakeNum(acc)
	}
	for _, n := range args {
		n, ok := n.(*Num)
		if !ok {
			panic("/ only takes number arguments")
		}

		acc *= n.n
	}
	return MakeNum(acc)
}

func DivPrim(o Obj, e *Env) Obj {
	acc := int64(1)
	args := listToSlice(Evlis(o, e))
	if len(args) == 0 {
		return MakeNum(acc)
	}
	for i, n := range args {
		n, ok := n.(*Num)
		if !ok {
			panic("/ only takes number arguments")
		}

		// first element is divident, following are divisors
		if i == 0 {
			acc = n.n
		} else {
			acc /= n.n
		}
	}
	return MakeNum(acc)
}

func ModuloPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 2 {
		panic("modulo takes 2 args")
	}

	v1, ok := args[0].(*Num)
	if !ok {
		panic("modulo only takes number arguments")
	}
	v2, ok := args[1].(*Num)
	if !ok {
		panic("modulo only takes number arguments")
	}

	return MakeNum(v1.n % v2.n)
}

func ExitPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) > 1 {
		panic("exit takes 1 or 0 argument")
	}
	if len(args) == 0 {
		os.Exit(0)
	}
	n, ok := args[0].(*Num)
	if !ok {
		panic("exit take a number for an argument")
	}
	os.Exit(int(n.n))
	return Nil
}
