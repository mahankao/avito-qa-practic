package utils

import "strings"

func MixCase(text string) string {
	letters := strings.Split(strings.ToLower(text), "")

	for i := range letters {
		if i%2 == 0 {
			letters[i] = strings.ToUpper(letters[i])
		}
	}

	return strings.Join(letters, "")
}
