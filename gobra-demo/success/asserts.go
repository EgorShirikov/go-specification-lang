package success

func AssertMul(x int) {
	/*@ assert x * 2 == x + x; */
	_ = x
}

func AssertSub(x int) {
	/*@ assert x - x == 0; */
	_ = x
}

func AssertNeq(x int) {
	/*@ assert x != x + 1; */
	_ = x
}

func AssertPositive(x int) {
	/*@ assert x >= 0 ==> x + 1 > 0; */
	_ = x
}

func AssertDouble(x int) {
	/*@ assert x >= 0 ==> x * 2 >= x; */
	_ = x
}
