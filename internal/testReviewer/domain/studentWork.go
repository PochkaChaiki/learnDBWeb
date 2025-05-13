package domain

// type StudentWork struct {
// 	Name  string
// 	Group string
// 	Tasks []TestTask
// 	DB    string
// }

type StudentWork struct {
	Name       string
	Group      string
	DB         string
	TotalGrade int
	Tasks      []Task
}

type Works []*StudentWork

// type WorkReview struct {
// 	Name       string
// 	Group      string
// 	DB         string
// 	TotalGrade int
// 	Tasks      []Task
// }

type ReviewedWorks []StudentWork
