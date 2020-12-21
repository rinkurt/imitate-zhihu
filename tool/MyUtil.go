package tool

import "strconv"

// Init for package tool
func init() {
	initConfig()
	initLogger()
	initDatabase()
	initRedis()
}

func If(cond bool, t interface{}, f interface{}) interface{} {
	if cond {
		return t
	}
	return f
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

