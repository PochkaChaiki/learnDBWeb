package domain

type Question struct {
	Id            int    `json:"question_id,omitempty" db:"question_id"`
	QuestionText  string `json:"question_text" db:"question_text"`
	CorrectAnswer string `json:"correct_answer" db:"correct_answer"`
	DBSampleId    int    `json:"dbsample_id" db:"db_sample_id"`
}
