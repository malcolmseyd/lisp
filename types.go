package main

import (
	"fmt"
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
var _ Obj = &Symbol{}
var _ Obj = &CloseParen{}
var _ Obj = &Pair{}
var _ Obj = &Num{}

// Symbol is an interned string
type Symbol struct {
	interned *string
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
	return &Symbol{interned: interned}
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

type Num struct {
	n int64
}

func (Num) Type() ObjType {
	return TypeNumber
}

func MakeNum(n int64) *Num {
	return &Num{n: n}
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

func (e *Env) Set(sym *Symbol, o Obj) {
	e.bindings[*sym] = o
}

func (e *Env) SetMut(sym *Symbol, o Obj) {
	if _, ok := e.bindings[*sym]; ok {
		e.bindings[*sym] = o
	}
	if e.parent != nil {
		e.parent.SetMut(sym, o)
	}
	panic(fmt.Sprintf("tried to set unbound variable %v", sym))
}

func (e *Env) Get(sym *Symbol) Obj {
	if o, ok := e.bindings[*sym]; ok {
		return o
	}
	if e.parent != nil {
		return e.Get(sym)
	}
	panic(fmt.Sprintf("tried to get unbound variable %v", sym))
}