package domain

type Task struct {
	ID                 int         `db:"id"`
	Name               *string     `db:"name"`
	Question           string      `db:"question"`
	Points             int         `db:"points"`
	GeneralFeedback    *string     `db:"general_feedback"`
	CorrectAnswerJSONB []byte      `db:"correct_answer_jsonb"`
	DBSamples          []DBSample  `db:"-"`
	XMLAnswers         []XMLAnswer `db:"-"`
	UserGroups         []string    `db:"-"`
}

type DBSample struct {
	ID          int     `db:"id"`
	Name        string  `db:"name"`
	Path        string  `db:"path"`
	Description *string `db:"description"`
	DBMS        []DBMS  `db:"-"`
}

type DBMS struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Test struct {
	ID         int        `db:"id"`
	Name       string     `db:"name"`
	Categories []Category `db:"-"`
}

type XMLAnswer struct {
	ID       int     `db:"id"`
	Text     string  `db:"text"`
	Feedback *string `db:"feedback"`
	Fraction float32 `db:"fraction"`
}

type UserGroup struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type User struct {
	ID       int      `db:"id"`
	Login    string   `db:"login"`
	Password string   `db:"password"`
	Groups   []string `db:"-"`
}

type Category struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Tasks []Task `db:"-"`
}

type ActiveDBInstance struct {
	ID               int      `db:"id"`
	ConnectionString string   `db:"connection_string"`
	DBMS             DBMS     `db:"-"`
	DBSample         DBSample `db:"-"`
}
