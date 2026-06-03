package failure

func AssertAdd(x, y int) {
	/*@ assert x + y >= x; */ // FAILS: y can be negative
	_ = x
	_ = y
}

func AssertAlwaysTrue(x int) {
	/*@ assert x == x + 1; */ // FAILS: always false
	_ = x
}
