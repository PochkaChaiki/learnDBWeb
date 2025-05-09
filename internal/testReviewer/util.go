package testReviewer

import (
	"regexp"
)

func RetrieveScripts(ans string) []string {
	re, err := regexp.Compile("(?i)(with[^;]*)?select(([\\'\\\"]\\s?;\\s?[\\'\\\"])|[^;])*")
	if err != nil {
		return nil
	}

	return re.FindAllString(ans, -1)
}

//
// // Saved this one function cause it can be useful in the future
//
// func RetrieveScript2(ans string) string {
// 	end := ";"
// 	semicolon := ";"
// 	lowerCasedAns := strings.ToLower(ans)

// 	selectIndex := strings.Index(lowerCasedAns, "select")
// 	withIndex := strings.Index(lowerCasedAns, "with")

// 	if selectIndex == -1 {
// 		return ""
// 	}
// 	startIndex := selectIndex

// 	if withIndex != -1 && withIndex < startIndex {
// 		startIndex = withIndex
// 	}

// 	if startIndex != 0 {
// 		switch symbolBefore := []rune(lowerCasedAns)[startIndex-1]; symbolBefore {
// 		case '"':
// 			end = "\""
// 		case '(':
// 			end = ")"
// 		case '[':
// 			end = "]"
// 		case '\'':
// 			end = "'"
// 		default:
// 		}
// 	}

// 	endIndex := strings.LastIndex(lowerCasedAns, end)
// 	if semicolonIndex := strings.LastIndex(lowerCasedAns, semicolon); semicolonIndex > endIndex {
// 		endIndex = semicolonIndex
// 	}

// 	if endIndex != -1 && endIndex > startIndex {
// 		if ans[endIndex-1:endIndex] == semicolon {
// 			return ans[startIndex : endIndex-1]
// 		}
// 		return ans[startIndex:endIndex]
// 	}

// 	return ans[startIndex:]
// }
