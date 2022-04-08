package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type ObjType uint8

const (
	// core types for the Lisp machine
	TypeSymbol = ObjType(iota)
	TypePair
	TypePrimitive
	TypeProcedure
	TypeMacro
	// parsing types
	TypeCloseParen
	// data types
	TypeNumber
)

// All Lisp objects must satisfy this interface
type Obj interface {
	Type() ObjType
}

var _ Obj = Primitive(nil)
var _ Obj = &Procedure{}
var _ Obj = &Macro{}
var _ Obj = &Symbol{}
var _ Obj = &CloseParen{}
var _ Obj = &Pair{}
var _ Obj = &Number{}

// Symbol is an interned string (except with Gensym)
type Symbol struct {
	s *string
}

func (s *Symbol) Type() ObjType {
	return TypeSymbol
}

func Intern(s string) *Symbol {
	interned, ok := symbols[s]
	if !ok {
		symbols[s] = &s
		interned = &s
	}
	return &Symbol{s: interned}
}

// names can collide since it's not interned
// the name is for debugging purposes only
var gensymCounter uint64 = 0

// generates an un-interned symbol
// for use inside of macros only please and thank you
func Gensym() *Symbol {
	s := "__GEN-" + strconv.FormatUint(gensymCounter, 36)
	gensymCounter++
	return &Symbol{s: &s}
}

func (s *Symbol) Equal(o Obj) bool {
	if sym, ok := o.(*Symbol); ok {
		return *sym == *s
	}
	return false
}

// Pair is a cons-cell
type Pair struct {
	Car Obj
	Cdr Obj
}

func (p *Pair) Type() ObjType {
	return TypePair
}

func Cons(car, cdr Obj) *Pair {
	return &Pair{Car: car, Cdr: cdr}
}

func Car(p *Pair) Obj {
	return p.Car
}

func Cdr(p *Pair) Obj {
	return p.Cdr
}

type CloseParen struct{}

func (CloseParen) Type() ObjType {
	return TypeCloseParen
}

type Primitive func(Obj, *Env) Obj

func (Primitive) Type() ObjType {
	return TypePrimitive
}

type Procedure struct {
	args     []Symbol
	body     Obj
	scope    *Env
	variadic *Symbol // nil if not variadic
}

func (Procedure) Type() ObjType {
	return TypeProcedure
}

func MakeProcedure(args []Symbol, body Obj, scope *Env) *Procedure {
	return &Procedure{args: args, body: body, scope: scope}
}

func MakeVariadicProcedure(args []Symbol, variadic Symbol, body Obj, scope *Env) *Procedure {
	return &Procedure{args: args, body: body, scope: scope, variadic: &variadic}
}

type Macro struct {
	args     []Symbol
	body     Obj
	scope    *Env
	variadic *Symbol // nil if not variadic
}

func (Macro) Type() ObjType {
	return TypeMacro
}

func MakeMacro(args []Symbol, body Obj, scope *Env) *Macro {
	return &Macro{args: args, body: body, scope: scope}
}

func MakeVariadicMacro(args []Symbol, variadic Symbol, body Obj, scope *Env) *Macro {
	return &Macro{args: args, body: body, scope: scope, variadic: &variadic}
}

type Number struct {
	n *big.Int
}

func (Number) Type() ObjType {
	return TypeNumber
}

func ParseNum(text []byte) *Number {
	n := big.NewInt(0)
	n.UnmarshalText(text)
	return &Number{n: n}
}

func MakeNum(n *big.Int) *Number {
	if n == nil {
		return &Number{n: big.NewInt(0)}
	}
	return &Number{n: n}
}

type Env struct {
	bindings map[Symbol]Obj
	parent   *Env
}

func MakeEnv(parent *Env) *Env {
	return &Env{
		bindings: map[Symbol]Obj{},
		parent:   parent,
	}
}

func (e *Env) Bind(sym *Symbol, o Obj) Obj {
	e.bindings[*sym] = o
	return o
}

func (e *Env) Set(sym *Symbol, o Obj) Obj {
	if old, ok := e.bindings[*sym]; ok {
		e.bindings[*sym] = o
		return old
	}
	if e.parent != nil {
		return e.parent.Set(sym, o)
	}
	panic(fmt.Sprintf("tried to set unbound variable %v", sym))
}

func (e *Env) Resolve(sym *Symbol) Obj {
	if o, ok := e.bindings[*sym]; ok {
		return o
	}
	if e.parent != nil {
		return e.parent.Resolve(sym)
	}
	panic(fmt.Sprintf("tried to get unbound variable %v", sym))
}

func (e *Env) String() string {
	s := strings.Builder{}
	s.WriteString("env ")
	s.WriteString(fmt.Sprintf("%p", e))
	s.WriteString(" {\n")
	for k, v := range e.bindings {
		s.WriteString(k.String())
		s.WriteString(": ")
		s.WriteString(v.(fmt.Stringer).String())
		s.WriteString(" {")
		s.WriteString(fmt.Sprintf("%p", v))
		s.WriteString("}")
		s.WriteString(",\n")
	}
	s.WriteString("}\n")
	return s.String()
}
