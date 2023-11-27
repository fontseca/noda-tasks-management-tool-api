package service

import "strings"

func doTrim(args ...*string) (trimmed int) {
	if 0 == len(args) {
		return 0
	}
	for _, str := range args {
		if nil != str {
			*str = strings.Trim(*str, " \a\b\f\n\r\t\v")
			trimmed++
		}
	}
	return trimmed
}
