package failure

/*
@
requires x >= 0;
ensures res == x + 2;  // WRONG: should be x + 1
assigns \nothing;
*/
func AddOneWrong(x int) (res int) {
	return x + 1
}

/*
@
requires b > a;
ensures res == b + a;  // WRONG: should be b - a
assigns \nothing;
*/
func DiffWrong(a, b int) (res int) {
	return b - a
}

/*
@
requires n >= 0;
ensures res == n * 3;  // WRONG: should be n * 2
assigns \nothing;
*/
func DoubleWrong(n int) (res int) {
	return n * 2
}
