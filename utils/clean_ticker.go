package utils

import (
	"strings"
)

func CleanTicker(ticker string) string {
	var sb strings.Builder
	for _, c := range ticker {
		if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) || c == 45 || c == 46 {
			sb.WriteRune(c)
		}
	}

	return strings.ToUpper(sb.String())
}
