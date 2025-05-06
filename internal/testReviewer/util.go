package testReviewer

import "strings"

func RetrieveScript(ans string) string {
	end := ";"
	lowerCasedAns := strings.ToLower(ans)

	selectIndex := strings.Index(lowerCasedAns, "select")
	withIndex := strings.Index(lowerCasedAns, "with")

	if selectIndex == -1 {
		return ""
	}
	startIndex := selectIndex

	if withIndex != -1 && withIndex < startIndex {
		startIndex = withIndex
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

	if endIndex := strings.LastIndex(lowerCasedAns, end); endIndex != -1 && endIndex > startIndex {
		return ans[startIndex:endIndex]
	}

	return ans[startIndex:]
}
