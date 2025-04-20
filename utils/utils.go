package utils

func If[T any](cond bool, tval, fval T) T {
	if cond {
		return tval
	}
	return fval
}

func GetRedisKey(namespace string, k string) string {
	return namespace + ":" + k
}
