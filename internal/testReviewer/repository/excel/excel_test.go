package excel

import (
	"learnDB/internal/testReviewer/config"
	"testing"
)

func TestExcelRead(t *testing.T) {
	staticPath := "/home/pochka/projects/learnDB/excelFiles/"
	config := config.MustLoadConfig(staticPath + "excelConfig.json")

	tr := &ExcelRepository{
		bookRead:  staticPath + "test.xlsx",
		sheetRead: "Sheet1",
	}

	res, err := tr.ReadFile(config)
	if err != nil {
		t.Errorf("excelReader read error: %v", err)
	}
	if res == nil {
		t.Errorf("got %v, want to nil", res)
	}
	if len(res) == 0 {
		t.Fatalf("len(res) == 0, res: %v", res)
	}
}

func TestExcelWrite(t *testing.T) {
	staticPath := "/home/pochka/projects/learnDB/excelFiles/"
	config := config.MustLoadConfig(staticPath + "excelConfig.json")

	er := &ExcelRepository{
		bookRead:   staticPath + "test.xlsx",
		sheetRead:  "Sheet1",
		bookWrite:  staticPath + "reviewedWorks.xlsx",
		sheetWrite: "Sheet1",
	}
	studentWorks, err := er.ReadFile(config)
	if err != nil {
		t.Fatalf("excelReader read error: %v", err)
	}

	// reviewedWorks := make(domain.ReviewedWorks, 0, len(studentWorks))
	// for _, work := range studentWorks {
	// 	taskReviews := make([]domain.TaskReview, 0, len(work.Tasks))
	// 	for _, task := range work.Tasks {
	// 		taskReviews = append(taskReviews, domain.TaskReview{
	// 			Task: task,
	// 			Review: domain.CheckResult{
	// 				Points: 3,
	// 				Error:  nil,
	// 			},
	// 		})
	// 	}
	// 	reviewedWorks = append(reviewedWorks, domain.WorkReview{
	// 		Name:       work.Name,
	// 		Group:      work.Group,
	// 		DB:         work.DB,
	// 		TotalGrade: 12,
	// 		Tasks:      taskReviews,
	// 	})
	// }

	err = er.WriteFile(studentWorks)
	if err != nil {
		t.Errorf("got error %v, want nil", err)
	}

}

func TestRetrieveScript(t *testing.T) {

	er := New("Sheet1")

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"with ;", "select * from query;", []string{"select * from query"}},
		{"without ;", "select * from query", []string{"select * from query"}},
		// {"square brackets", `[select * from query]`, []string{"select * from query"}},
		// {"double quotemarks", `"select * from query"`, []string{"select * from query"}},
		// {"single quotemarks", `'select * from query'`, []string{"select * from query"}},
		// {"parentheses", `(select * from query)`, []string{"select * from query"}},
		{"answer ... parentheses", `501 (SELECT * FROM query;)`, []string{"SELECT * FROM query"}},
		{"answer ... broken parentheses", `10 (SELECT MIN(age) FROM olympics.games_competito r;`, []string{"SELECT MIN(age) FROM olympics.games_competito r"}},
		{"nested query", `select * from (select id from query) s limit 1;`, []string{"select * from (select id from query) s limit 1"}},
		{
			"answer ... query",
			`Shooting; Athletics; Swimming SELECT STRING_AGG(sport_name, '; ') AS "Топ-3 вида спорта" FROM (SELECT olympics.sport.sport_name FROM olympics.sport JOIN olympics.event e ON olympics.sport.id = e.sport_id GROUP BY olympics.sport.sport_name ORDER BY COUNT(e.id) DESC LIMIT 3 );`,
			[]string{`SELECT STRING_AGG(sport_name, '; ') AS "Топ-3 вида спорта" FROM (SELECT olympics.sport.sport_name FROM olympics.sport JOIN olympics.event e ON olympics.sport.id = e.sport_id GROUP BY olympics.sport.sport_name ORDER BY COUNT(e.id) DESC LIMIT 3 )`},
		},
		{
			"postgres query starts with 'with'",
			`WITH top_sports AS (SELECT s.id AS sport_id, s.sport_name, COUNT(DISTINCT e.id) AS event_count FROM sport s JOIN event e ON s.id = e.sport_id GROUP BY s.id, s.sport_name ORDER BY event_count DESC LIMIT 3 ) SELECT ts.sport_name, ts.event_count AS medal_events_count FROM top_sports ts UNION ALL SELECT 'TOTAL' AS sport_name, SUM(ts.event_count) AS medal_events_count FROM  top_sports ts;`,

			[]string{`WITH top_sports AS (SELECT s.id AS sport_id, s.sport_name, COUNT(DISTINCT e.id) AS event_count FROM sport s JOIN event e ON s.id = e.sport_id GROUP BY s.id, s.sport_name ORDER BY event_count DESC LIMIT 3 ) SELECT ts.sport_name, ts.event_count AS medal_events_count FROM top_sports ts UNION ALL SELECT 'TOTAL' AS sport_name, SUM(ts.event_count) AS medal_events_count FROM  top_sports ts`},
		},
		{"empty answer", "", nil},
		{"no script with \"-\"", "-", nil},
		{"no script with some answer", "2025", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := er.retrieveScripts(test.input)
			if got == nil && test.want == nil {
				return
			}
			if got == nil || test.want == nil {
				t.Errorf("got %q, \nwant %q", got, test.want)
			}
			if len(got) != len(test.want) {
				t.Errorf("got %q, \nwant %q", got, test.want)
			}
			for i := range got {
				if got[i] != test.want[i] {
					t.Errorf("got %q, \nwant %q", got, test.want)
				}
			}
		})

	}
}
