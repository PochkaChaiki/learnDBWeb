package domain

type WorkReview []CheckResult

type StudentWork struct {
	Name  string
	Group string
	Tasks []TestTask
	DB    string
}

type ReviewedStudentWork struct {
	TotalGrade int
	StudentWork
	WorkReview WorkReview
}

type ReviewedWorks []ReviewedStudentWork
