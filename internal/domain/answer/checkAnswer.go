package answer

import (
	"fmt"
	"math"
	"strings"
)

func StrictCheck(data []any, corrAns CorrectAnswer) (int, bool) {
	if len(corrAns.Values) != len(data) {
		return 0, false
	}

	ansIsCorrect := true
	// If every part of student's answer is present at the correct answer than the first one is correct
	for _, ans := range data {
		var ansToCheck string
		switch ans := ans.(type) {
		case string:
			ansToCheck = ans
		case []uint8:
			ansToCheck = string(ans)
		default:
			ansToCheck = fmt.Sprint(ans)
		}

		partIsCorrect := false
		for _, corrAnsPart := range corrAns.Values {
			partIsCorrect = partIsCorrect || strings.EqualFold(ansToCheck, corrAnsPart)
			if partIsCorrect {
				break
			}
		}
		ansIsCorrect = ansIsCorrect && partIsCorrect
	}

	if ansIsCorrect {
		return corrAns.Points, true
	}
	return 0, false
}

func LightCheck(data []any, corrAns CorrectAnswer) (int, bool) {
	coef := float64(len(corrAns.Values)) / float64(len(data))
	if coef > 1 {
		return 0, false
	}

	ansIsCorrect := true
	for _, corrAnsPart := range corrAns.Values {
		partIsCorrect := false
		for _, ans := range data {
			var ansToCheck string
			switch ans := ans.(type) {
			case string:
				ansToCheck = ans
			case []uint8:
				ansToCheck = string(ans)
			default:
				ansToCheck = fmt.Sprint(ans)
			}
			partIsCorrect = partIsCorrect || strings.EqualFold(ansToCheck, corrAnsPart)
			if partIsCorrect {
				break
			}
		}
		ansIsCorrect = ansIsCorrect && partIsCorrect
	}
	if ansIsCorrect {
		return int(math.Round(float64(corrAns.Points) * coef)), true
	}
	return 0, false
}
