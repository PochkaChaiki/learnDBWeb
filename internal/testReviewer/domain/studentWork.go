package domain

type StudentWork struct {
	Name  string
	Group string
	Tasks []TestTask
}

type WorkReview []CheckResult

type ReviewedStudentWork struct {
	TotalGrade int
	StudentWork
	WorkReview WorkReview
}

type ReviewedWorks []ReviewedStudentWork
