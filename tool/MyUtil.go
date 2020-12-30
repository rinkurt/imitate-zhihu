package tool

import (
	"github.com/thinkeridea/go-extend/exunicode/exutf8"
	"strconv"
	"strings"
)

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

func StrToDefaultInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StrToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ParseCursor(cursor string) ([]int64, error) {
	if cursor == "" {
		return []int64{-1, -1}, nil
	}
	cursors := strings.Split(cursor, ",")
	ret := make([]int64, len(cursors))
	for i, v := range cursors {
		el, err := StrToInt64(v)
		if err != nil {
			return nil, err
		}
		ret[i] = el
	}
	return ret, nil
}

func CutContent(content *string, length int) {
	if len(*content) > length {
		*content = exutf8.RuneSubString(*content, 0, length) + "..."
	}
}
