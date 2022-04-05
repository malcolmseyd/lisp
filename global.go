package main

var symbols = map[string]*string{}

var (
	Nil      = Intern("nil")
	True     = Intern("#t")
	Dot      = Intern(".")
	QuoteSym = Intern("quote")
)

func BindGlobals(e *Env) {
	prims := map[string]func(Obj, *Env) Obj{
		"lambda":   LambdaPrim,
		"cons":     ConsPrim,
		"car":      CarPrim,
		"cdr":      CdrPrim,
		"define":   DefinePrim,
		"defmacro": DefMacroPrim,
		"set!":     SetPrim,
		"set-car!": SetCarPrim,
		"set-cdr!": SetCdrPrim,
		"if":       IfPrim,
		"eq?":      EqPrim,
		"=":        EqPrim,
		"<":        LessPrim,
		"quote":    QuotePrim,
		"eval":     EvalPrim,
		"apply":    ApplyPrim,
		"+":        AddPrim,
		"-":        SubPrim,
		"*":        MulPrim,
		"/":        DivPrim,
		"modulo":   ModuloPrim,
		"exit":     ExitPrim,
		"print":    PrintPrim,
	}

	for name, f := range prims {
		e.Bind(Intern(name), Primitive(f))
	}

	e.Bind(Intern("nil"), Nil)
	e.Bind(Intern("#t"), True)
}
