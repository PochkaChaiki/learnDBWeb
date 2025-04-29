package testReviewer

import "strings"

func RetrieveScript(ans string) string {
	start := "select"
	end := ";"
	lowerCasedAns := strings.ToLower(ans)

	startIndex := strings.Index(lowerCasedAns, start)
	if startIndex == -1 {
		return ""
	}

	if startIndex != 0 {
		switch symbolBefore := []rune(lowerCasedAns)[startIndex-1]; symbolBefore {
		case '"':
			end = "\""
		case '(':
			end = ")"
		case '[':
			end = "]"
		case '\'':
			end = "'"
		default:
		}
	}

	if endIndex := strings.Index(lowerCasedAns[startIndex:], end); endIndex != -1 {
		return lowerCasedAns[startIndex : endIndex+1]
	}

	return lowerCasedAns[startIndex:]
}
