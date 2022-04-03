package main

func boolToLisp(b bool) Obj {
	if b {
		return True
	}
	return Nil
}
