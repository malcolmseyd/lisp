package main

var symbols = map[string]*string{}

// TODO: top level environment frame

var (
	Nil = Intern("nil")
)

func BindGlobals(e *Env) {
	bindPrim("cons", ConsPrim, e)
	bindPrim("quote", QuotePrim, e)
	bindPrim("eval", EvalPrim, e)
	e.Bind(Intern("nil"), Nil)
}

func bindPrim(name string, f func(Obj, *Env) Obj, e *Env) {
	e.Bind(Intern(name), Primitive(f))
}
