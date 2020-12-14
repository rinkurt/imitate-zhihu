package tool

import "strconv"

const (
	OrderByTime = iota
	OrderByHeat
)

func If(cond bool, t interface{}, f interface{}) interface{} {
	if cond {
		return t
	}
	return f
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

