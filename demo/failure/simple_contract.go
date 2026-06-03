package failure

/*
@
requires x >= 0;
ensures res == x + 2;  // WRONG: should be x + 1
assigns \nothing;
*/
func AddOneWrongSimple(x int) (res int) {
	return x + 1
}
