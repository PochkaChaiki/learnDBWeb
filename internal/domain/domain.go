package domain

type Task struct {
	ID                 int     `db:"id"`
	Name               *string `db:"name"`
	Question           string  `db:"question"`
	Points             int     `db:"points"`
	GeneralFeedback    *string `db:"general_feedback"`
	CorrectAnswerJSONB *string `db:"correct_answer_jsonb"`
	DBSamples          []DBSample
	Answers            []Answer
	UserGroups         []string
}

type DBSample struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Path        string  `db:"path"`
	Description *string `db:"description"`
	DBMS        *string `db:"dbms"`
}

type DBMS struct {
	Name string `db:"name"`
}

type Test struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Categories []Category
}

type Answer struct {
	ID       int     `db:"id"`
	Text     string  `db:"text"`
	Feedback *string `db:"feedback"`
}

type UserGroup struct {
	Name string `db:"name"`
}

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Groups   []string
}

type Category struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Tasks []Task
}
