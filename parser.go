package main

import (
	"io"
	"strings"
	"unicode"
)

/*
The parser will use an LL recursive descent parsing strategy
*/

func readRune(s io.RuneScanner) rune {
	r, _, e := s.ReadRune()
	if e != nil {
		panic(e.Error())
	}
	return r
}

func peekRune(s io.RuneScanner) rune {
	r, _, e := s.ReadRune()
	if e != nil {
		panic(e.Error())
	}
	s.UnreadRune()
	return r
}

func Read(s io.RuneScanner) Obj {
	ReadSpace(s)
	for ReadComment(s) {
		ReadSpace(s)
	}
	if o := ReadList(s); o != nil {
		return o
	}
	if o := ReadCloseParen(s); o != nil {
		return o
	}
	if o := ReadNum(s); o != nil {
		return o
	}
	if o := ReadQuote(s); o != nil {
		return o
	}
	// this should be at the bottom since it's so permissive
	if o := ReadSym(s); o != nil {
		return o
	}
	panic("bug: unknown syntax encountered while reading")
}

func ReadSpace(s io.RuneScanner) {
	for unicode.IsSpace(peekRune(s)) {
		readRune(s)
	}
}

func ReadComment(s io.RuneScanner) bool {
	if peekRune(s) == ';' {
		for readRune(s) != '\n' {
			// consume until newline
		}
		return true
	}
	return false
}

func ReadQuote(s io.RuneScanner) Obj {
	r := peekRune(s)
	if r != '\'' {
		return nil
	}
	readRune(s)
	return Cons(Intern("quote"), Cons(Read(s), Nil))
}

const symbolChars = "!#$%&*+,-./@:<=>?^_"

func isSymRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || strings.ContainsRune(symbolChars, r)
}

func ReadSym(s io.RuneScanner) *Symbol {
	b := strings.Builder{}
	for r := peekRune(s); isSymRune(r); r = peekRune(s) {
		readRune(s)
		b.WriteRune(r)
	}
	if b.Len() == 0 {
		return nil
	}
	return Intern(b.String())
}

func isNumRune(r rune) bool {
	return r >= '0' && r <= '9'
}

func parseDigit(r rune) int64 {
	return int64(r - '0')
}

func ReadNum(s io.RuneScanner) *Num {
	r := peekRune(s)
	result := int64(0)
	if isNumRune(r) {
		result = result*10 + parseDigit(r)
		readRune(s)
		r = peekRune(s)
	} else {
		return nil
	}
	for isNumRune(r) {
		result = result*10 + parseDigit(r)
		readRune(s)
		r = peekRune(s)
	}
	return MakeNum(result)
}

func ReadList(s io.RuneScanner) Obj {
	open := peekRune(s)
	if open != '(' {
		return nil
	}
	readRune(s)

	start := Obj(Nil)
	prev := start
Outer:
	for {
		switch curr := Read(s).(type) {
		case *CloseParen:
			break Outer
		default:
			if dot, ok := curr.(*Symbol); ok && *dot == *Dot {
				curr = Read(s)
				if Nil.Equal(start) {
					start = curr
				}
				if prevPair, ok := prev.(*Pair); ok && !Nil.Equal(prev) {
					prevPair.Cdr = curr
				}
				if Read(s).Type() != TypeCloseParen {
					panic("missing close paren after .")
				}
				break Outer
			}
			curr = Cons(curr, Nil)
			if Nil.Equal(start) {
				start = curr
			}
			if prevPair, ok := prev.(*Pair); ok && !Nil.Equal(prev) {
				prevPair.Cdr = curr
			}
			prev = curr
		}
	}
	return start
}

func ReadCloseParen(s io.RuneScanner) *CloseParen {
	if peekRune(s) == ')' {
		readRune(s)
		return &CloseParen{}
	}
	return nil
}
