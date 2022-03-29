package main

import "io"

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

// TODO: write parser
