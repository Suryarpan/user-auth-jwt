package utils

func If[T any](cond bool, tval, fval T) T {
	if cond {
		return tval
	}
	return fval
}
