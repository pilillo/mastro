package strings

import "strings"

func SplitAndTrim(input string, separator string) []string {
	vals := strings.Split(input, separator)
	for i := range vals {
		vals[i] = strings.TrimSpace(vals[i])
	}
	return vals
}
