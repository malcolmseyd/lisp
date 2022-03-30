package main

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

// Symbol is an interned string
type Symbol struct {
	interned *string
}

var _ Obj = &Symbol{}

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

var _ Obj = &Pair{}

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

var _ Obj = &CloseParen{}

func (CloseParen) Type() ObjType {
	return TypeCloseParen
}
