package main

import "fmt"

func boolToLisp(b bool) Obj {
	if b {
		return True
	}
	return Nil
}

func sliceToList(slice []Obj) Obj {
	if len(slice) == 0 {
		return Nil
	}
	head := Cons(slice[0], Nil)
	tail := head
	for _, curr := range slice[1:] {
		next := Cons(curr, Nil)
		tail.Cdr = next
		tail = next
	}
	return head
}

func listToSlice(o Obj) []Obj {
	slice := make([]Obj, 0)
	for !Nil.Equal(o) {
		pair, ok := o.(*Pair)
		if !ok {
			panic(fmt.Sprintf("expected list, got %v", o))
		}
		slice = append(slice, Car(pair))
		o = Cdr(pair)
	}
	return slice
}

// for a list potentially not ending in Nil, like in a variadic function
func improperListToSlice(o Obj) ([]Obj, Obj) {
	slice := make([]Obj, 0)
	for !Nil.Equal(o) {
		pair, ok := o.(*Pair)
		if !ok {
			return slice, o
		}
		slice = append(slice, Car(pair))
		o = Cdr(pair)
	}
	return slice, nil
}
