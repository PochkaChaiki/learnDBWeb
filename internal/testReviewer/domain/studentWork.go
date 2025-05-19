package domain

// Struct to pass it to the testReviewer
type Work struct {
	DB         string
	TotalGrade int
	Units      []*CheckUnit
}

type Works []*Work
