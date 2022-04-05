package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

/*
The parser will use an LL recursive descent parsing strategy
*/

func readRune(s *bufio.Reader) rune {
	r, _, e := s.ReadRune()
	if e != nil {
		panic(e.Error())
	}
	return r
}

func peekRune(s *bufio.Reader) rune {
	r, _, e := s.ReadRune()
	if e != nil {
		panic(e.Error())
	}
	s.UnreadRune()
	return r
}

func peekN(s *bufio.Reader, n int) []byte {
	b, err := s.Peek(n)
	if err != nil {
		panic(fmt.Sprintf("error peeking %v bytes: %v", n, err))
	}
	return b
}

func consumeN(s *bufio.Reader, n int) {
	_, err := io.CopyN(io.Discard, s, int64(n))
	if err != nil {
		panic(fmt.Sprintf("error peeking %v bytes: %v", n, err))
	}
}

func Read(s *bufio.Reader) Obj {
	// things to ignore
	ReadSpace(s)
	for ReadComment(s) {
		ReadSpace(s)
	}

	readers := []func(*bufio.Reader) Obj{
		ReadList,
		ReadCloseParen,
		ReadNum,
		ReadQuote,
		ReadQuasiquote,
		ReadUnquoteSplicing, // must be above ReadUnquote
		ReadUnquote,
		ReadSym, // this should be at the bottom since it's so permissive
	}

	for _, reader := range readers {
		if o := reader(s); o != nil {
			return o
		}
	}
	panic("bug: unknown syntax encountered while reading")
}

func ReadSpace(s *bufio.Reader) {
	for unicode.IsSpace(peekRune(s)) {
		readRune(s)
	}
}

func ReadComment(s *bufio.Reader) bool {
	if peekRune(s) == ';' {
		for readRune(s) != '\n' {
			// consume until newline
		}
		return true
	}
	return false
}

func ReadQuote(s *bufio.Reader) Obj {
	r := peekRune(s)
	if r != '\'' {
		return nil
	}
	readRune(s)
	return Cons(QuoteSym, Cons(Read(s), Nil))
}

func ReadQuasiquote(s *bufio.Reader) Obj {
	r := peekRune(s)
	if r != '`' {
		return nil
	}
	readRune(s)
	return Cons(QuasiquoteSym, Cons(Read(s), Nil))
}

func ReadUnquoteSplicing(s *bufio.Reader) Obj {
	b := peekN(s, 2)
	if string(b) != ",@" {
		return nil
	}
	consumeN(s, 2)
	return Cons(UnquoteSplicingSym, Cons(Read(s), Nil))
}

func ReadUnquote(s *bufio.Reader) Obj {
	r := peekRune(s)
	if r != ',' {
		return nil
	}
	readRune(s)
	return Cons(UnquoteSym, Cons(Read(s), Nil))
}

const symbolChars = "!#$%&*+-./@:<=>?^_"

func isSymRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || strings.ContainsRune(symbolChars, r)
}

func ReadSym(s *bufio.Reader) Obj {
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

func ReadNum(s *bufio.Reader) Obj {
	r := peekRune(s)
	b := bytes.Buffer{}
	for isNumRune(r) {
		b.WriteRune(r)
		readRune(s)
		r = peekRune(s)
	}
	if b.Len() == 0 {
		return nil
	}
	return ParseNum(b.Bytes())
}

func ReadList(s *bufio.Reader) Obj {
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

func ReadCloseParen(s *bufio.Reader) Obj {
	if peekRune(s) == ')' {
		readRune(s)
		return &CloseParen{}
	}
	return nil
}
