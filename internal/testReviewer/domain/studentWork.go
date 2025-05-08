package domain

type StudentWork struct {
	Name  string
	Group string
	Tasks []TestTask
	DB    string
}

type WorkReview struct {
	Name       string
	Group      string
	DB         string
	TotalGrade int
	Tasks      []TaskReview
}

type ReviewedWorks []WorkReview
