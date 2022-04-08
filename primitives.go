package main

import (
	"log"
	"math/big"
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

func DefMacroPrim(o Obj, e *Env) Obj {
	formArgs := listToSlice(o)
	if len(formArgs) < 3 {
		panic("defmacro takes at least 3 arguments")
	}

	name, ok := formArgs[0].(*Symbol)
	if !ok {
		panic("name must be a symbol")
	}

	args, variadic := improperListToSlice(formArgs[1])

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

	body := sliceToList(formArgs[2:])

	proc := Obj(nil)
	if variadicSym != nil {
		proc = MakeVariadicMacro(argsSyms, *variadicSym, body, e)
	} else {
		proc = MakeMacro(argsSyms, body, e)
	}
	return e.Bind(name, proc)
}

// only goes one layer deep
// if i want recursion i basically need another specialized Eval
func MacroExpandPrim(o Obj, e *Env) Obj {
	args := listToSlice(o)
	if len(args) != 1 {
		panic("macroexpand takes 1 argument")
	}
	expr := listToSlice(args[0])
	if len(expr) < 1 {
		panic("macroexpand must be used on a macro expression")
	}
	macro, ok := Eval(expr[0], e).(*Macro)
	if !ok {
		panic("macroexpand must be used on a macro expression")
	}
	rest := sliceToList(expr[1:])

	return ApplyMacro(macro, rest, e)
}

func GensymPrim(o Obj, e *Env) Obj {
	if !Nil.Equal(o) {
		panic("gensym takes no args")
	}
	return Gensym()
}

func IsSymbolPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("symbol? takes 1 argument")
	}
	switch args[0].(type) {
	case *Symbol:
		return True
	default:
		return Nil
	}
}

func IsPairPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("pair? takes 1 argument")
	}
	switch args[0].(type) {
	case *Pair:
		return True
	default:
		return Nil
	}
}

func IsPrimitivePrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("primitive? takes 1 argument")
	}
	switch args[0].(type) {
	case Primitive:
		return True
	default:
		return Nil
	}
}

func IsProcedurePrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("procedure? takes 1 argument")
	}
	switch args[0].(type) {
	case *Procedure:
		return True
	default:
		return Nil
	}
}

func IsMacroPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("macro? takes 1 argument")
	}
	switch args[0].(type) {
	case *Macro:
		return True
	default:
		return Nil
	}
}

func IsNumberPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("number? takes 1 argument")
	}
	switch args[0].(type) {
	case *Number:
		return True
	default:
		return Nil
	}
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
	case *Number:
		return boolToLisp(v1.n.Cmp(v2.(*Number).n) == 0)
	default:
		// reference equality for misc
		return boolToLisp(v1 == v2)
	}
}

func LessPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))

	v1, ok := args[0].(*Number)
	if !ok {
		panic("args should be numbers")
	}

	v2, ok := args[1].(*Number)
	if !ok {
		panic("args should be numbers")
	}

	return boolToLisp(v1.n.Cmp(v2.n) == -1)
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
	if len(args) == 2 {
		test := Eval(args[0], e)
		expr1 := args[1]
		if !Nil.Equal(test) {
			return Eval(expr1, e)
		}
		return Nil
	} else if len(args) == 3 {
		test := Eval(args[0], e)
		expr1 := args[1]
		expr2 := args[2]
		if !Nil.Equal(test) {
			return Eval(expr1, e)
		}
		return Eval(expr2, e)
	} else {
		panic("if takes 2 or 3 arguments")
	}
}

func CondPrim(o Obj, e *Env) Obj {
	args := listToSlice(o)
	if len(args) < 1 {
		panic("cond takes at least 1 argument")
	}
	for _, cases := range args {
		cases := listToSlice(cases)
		if len(cases) != 2 {
			panic("each cond case should have a predicate and a body")
		}
		pred := cases[0]
		body := cases[1]
		if Else.Equal(pred) {
			return Eval(body, e)
		}
		if !Nil.Equal(Eval(pred, e)) {
			return Eval(body, e)
		}
	}
	return Nil
}

func QuotePrim(o Obj, _ *Env) Obj {
	args := listToSlice(o)
	if len(args) != 1 {
		panic("quote takes 1 argument")
	}
	return args[0]
}

func QuasiquotePrim(o Obj, e *Env) Obj {
	args := listToSlice(o)
	if len(args) != 1 {
		panic("quasiquote takes 1 argument")
	}
	return Quasiquote(args[0], e)
}

func Quasiquote(o Obj, e *Env) Obj {
	pair, ok := o.(*Pair)
	if !ok {
		if UnquoteSym.Equal(o) {
			// since we can't catch it earlier
			panic("unquote takes 2 args")
		}
		return o
	}
	elems := listToSlice(pair)
	if len(elems) == 2 && UnquoteSym.Equal(elems[0]) {
		return Eval(elems[1], e)
	}
	out := make([]Obj, 0, len(elems))
	for _, elem := range elems {
		if isUnquoteSplicing(elem) {
			splicing := listToSlice(elem)
			if len(splicing) != 2 {
				panic("unquote-splicing takes 2 args")
			}
			out = append(out, listToSlice(Eval(splicing[1], e))...)
		} else {
			out = append(out, Quasiquote(elem, e))
		}
	}
	return sliceToList(out)
}

func isUnquoteSplicing(o Obj) bool {
	pair, ok := o.(*Pair)
	if !ok {
		return false
	}
	return UnquoteSplicingSym.Equal(Car(pair))
}

// TODO gensym
// make a constructor in types.go
// should be un-interned, only used in macros

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
	acc := big.NewInt(0)
	args := listToSlice(Evlis(o, e))
	if len(args) == 0 {
		return MakeNum(acc)
	}
	for _, n := range args {
		n, ok := n.(*Number)
		if !ok {
			panic("+ only takes number arguments")
		}

		acc.Add(acc, n.n)
	}
	return MakeNum(acc)
}

func SubPrim(o Obj, e *Env) Obj {
	acc := big.NewInt(0)
	args := listToSlice(Evlis(o, e))
	if len(args) == 0 {
		return MakeNum(acc)
	}
	for i, n := range args {
		n, ok := n.(*Number)
		if !ok {
			panic("- only takes number arguments")
		}

		// first element is minuend, following are subtrahend (i googled this lol)
		if i == 0 {
			acc.SetBytes(n.n.Bytes())
		} else {
			acc.Sub(acc, n.n)
		}
	}
	// special case: unary minus is negation
	if len(args) == 1 {
		return MakeNum(acc.Neg(acc))
	}
	return MakeNum(acc)
}

func MulPrim(o Obj, e *Env) Obj {
	acc := big.NewInt(1)
	args := listToSlice(Evlis(o, e))
	if len(args) == 0 {
		return MakeNum(acc)
	}
	for _, n := range args {
		n, ok := n.(*Number)
		if !ok {
			panic("* only takes number arguments")
		}

		acc.Mul(acc, n.n)
	}
	return MakeNum(acc)
}

func DivPrim(o Obj, e *Env) Obj {
	acc := big.NewInt(1)
	args := listToSlice(Evlis(o, e))
	if len(args) == 0 {
		return MakeNum(acc)
	}
	for i, n := range args {
		n, ok := n.(*Number)
		if !ok {
			panic("/ only takes number arguments")
		}

		// first element is divident, following are divisors
		if i == 0 {
			acc.SetBytes(n.n.Bytes())
		} else {
			acc.Div(acc, n.n)
		}
	}
	return MakeNum(acc)
}

func ModuloPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 2 {
		panic("modulo takes 2 args")
	}

	v1, ok := args[0].(*Number)
	if !ok {
		panic("modulo only takes number arguments")
	}
	v2, ok := args[1].(*Number)
	if !ok {
		panic("modulo only takes number arguments")
	}

	// mod MUTATES, so make a new bigint
	result := big.NewInt(0).Mod(v1.n, v2.n)
	return MakeNum(result)
}

func ExitPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) > 1 {
		panic("exit takes 1 or 0 arguments")
	}
	if len(args) == 0 {
		os.Exit(0)
	}
	n, ok := args[0].(*Number)
	if !ok {
		panic("exit take a number for an argument")
	}
	os.Exit(int(n.n.Int64()))
	return Nil
}

func PrintPrim(o Obj, e *Env) Obj {
	args := listToSlice(Evlis(o, e))
	if len(args) != 1 {
		panic("print takes 1 argument")
	}
	Print(args[0])
	return Nil
}

func PrintEnvPrim(o Obj, e *Env) Obj {
	log.Printf("%v\n", e.String())
	return Nil
}
