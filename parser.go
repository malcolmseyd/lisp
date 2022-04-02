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
	r, _, _ := s.ReadRune()
	return r
}

func peekRune(s io.RuneScanner) rune {
	r, _, _ := s.ReadRune()
	s.UnreadRune()
	return r
}

func Read(s io.RuneScanner) Obj {
	ReadSpace(s)
	if o := ReadList(s); o != nil {
		return o
	}
	if o := ReadCloseParen(s); o != nil {
		return o
	}
	if o := ReadNum(s); o != nil {
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

const symbolChars = "!#$%&*+,-./@:;<=>?^_"

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
	neg := false
	result := int64(0)
	if r == '-' {
		neg = true
		readRune(s)
		r = peekRune(s)
	} else if isNumRune(r) {
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
	if neg {
		result *= -1
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
	for {
		switch curr := Read(s).(type) {
		case *CloseParen:
			return start
		default:
			curr = Cons(curr, Nil)
			if start == Nil {
				start = curr
			}
			if prevPair, ok := prev.(*Pair); ok && prev != Nil {
				prevPair.Cdr = curr
			}
			prev = curr
		}
	}
}

func ReadCloseParen(s io.RuneScanner) *CloseParen {
	if peekRune(s) == ')' {
		readRune(s)
		return &CloseParen{}
	}
	return nil
}
