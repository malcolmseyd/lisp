package main

var symbols = map[string]*string{}

var (
	Nil = Intern("nil")
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
		"quote":  QuotePrim,
		"eval":   EvalPrim,
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
}
