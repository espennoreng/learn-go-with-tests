package iteration

import "strings"

func Repeat(char string, times int) string {
	var rep strings.Builder
	for range times {
		rep.WriteString(char)
	}
	return rep.String()
}
