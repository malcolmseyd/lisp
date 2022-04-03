package main

var symbols = map[string]*string{}

var (
	Nil  = Intern("nil")
	True = Intern("#t")
	Dot  = Intern(".")
)

func BindGlobals(e *Env) {
	prims := map[string]func(Obj, *Env) Obj{
		"lambda": LambdaPrim,
		"cons":   ConsPrim,
		"car":    CarPrim,
		"cdr":    CdrPrim,
		"define": DefinePrim,
		"set!":   SetPrim,
		"if":     IfPrim,
		"eq?":    EqPrim,
		"=":      EqPrim,
		"<":      LessPrim,
		"quote":  QuotePrim,
		"eval":   EvalPrim,
		"apply":  ApplyPrim,
		"+":      AddPrim,
		"-":      SubPrim,
		"*":      MulPrim,
		"/":      DivPrim,
		"exit":   ExitPrim,
	}

	for name, f := range prims {
		e.Bind(Intern(name), Primitive(f))
	}

	e.Bind(Intern("nil"), Nil)
	e.Bind(Intern("#t"), True)
}
