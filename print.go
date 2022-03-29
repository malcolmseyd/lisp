package main

import (
	"fmt"
	"log"
	"strings"
)

func (s *Symbol) String() string {
	return *s.interned
}

func (p *Pair) String() string {
	b := strings.Builder{}
	b.WriteByte('(')

	firstElem := true
	curr := Obj(p)
	for curr.Type() == TypePair {
		pair := curr.(*Pair)
		if !firstElem {
			b.WriteByte(' ')
		} else {
			firstElem = false
		}

		b.WriteString(mustStringer(Car(pair)).String())
		curr = Cdr(pair)
	}
	if curr != Nil {
		b.WriteString(" . ")
		b.WriteString(mustStringer(curr).String())
	}
	b.WriteByte(')')
	return b.String()
}

var _ fmt.Stringer = &Symbol{}
var _ fmt.Stringer = &Pair{}

// see above for a list of supported types
func Print(o Obj) {
	fmt.Println(mustStringer(o))
}

// exit on any bugs, all user exposed types should be printable
func mustStringer(o Obj) fmt.Stringer {
	s, ok := o.(fmt.Stringer)
	if !ok {
		log.Fatalf("bug: cannot string type %T\n", o)
	}
	return s
}
