package testReviewer

import "testing"

func TestRetrieveScript(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"with ;", "select * from query;", "select * from query"},
		{"without ;", "select * from query", "select * from query"},
		{"double quotemarks", `"select * from query"`, "select * from query"},
		{"single quotemarks", `'select * from query'`, "select * from query"},
		{"parentheses", `(select * from query)`, "select * from query"},
		{"answer ... parentheses", `501 (SELECT * FROM query;)`, "SELECT * FROM query"},
		{"answer ... broken parentheses", `10 (SELECT MIN(age) FROM olympics.games_competito r;`, "SELECT MIN(age) FROM olympics.games_competito r"},
		{"square brackets", `[select * from query]`, "select * from query"},
		{"nested query", `select * from (select id from query) s limit 1;`, "select * from (select id from query) s limit 1"},
		{
			"answer ... query",
			`Shooting; Athletics; Swimming SELECT STRING_AGG(sport_name, '; ') AS "Топ-3 вида спорта" FROM (
     SELECT
          olympics.sport.sport_name
     FROM
          olympics.sport
     JOIN
          olympics.event e ON olympics.sport.id = e.sport_id
     GROUP BY
          olympics.sport.sport_name
     ORDER BY
          COUNT(e.id) DESC
     LIMIT 3 );`,
			`SELECT STRING_AGG(sport_name, '; ') AS "Топ-3 вида спорта" FROM (
     SELECT
          olympics.sport.sport_name
     FROM
          olympics.sport
     JOIN
          olympics.event e ON olympics.sport.id = e.sport_id
     GROUP BY
          olympics.sport.sport_name
     ORDER BY
          COUNT(e.id) DESC
     LIMIT 3 )`},
		{
			"postgres query starts with 'with'",
			`WITH top_sports AS (SELECT
        s.id AS sport_id,
        s.sport_name,
        COUNT(DISTINCT e.id) AS event_count
    FROM sport s
    JOIN event e ON s.id = e.sport_id
    GROUP BY s.id, s.sport_name
    ORDER BY event_count DESC
    LIMIT 3 )
 SELECT ts.sport_name, ts.event_count AS medal_events_count FROM top_sports ts UNION ALL SELECT 'TOTAL' AS sport_name, SUM(ts.event_count) AS medal_events_count FROM  top_sports ts;`,

			`WITH top_sports AS (SELECT
        s.id AS sport_id,
        s.sport_name,
        COUNT(DISTINCT e.id) AS event_count
    FROM sport s
    JOIN event e ON s.id = e.sport_id
    GROUP BY s.id, s.sport_name
    ORDER BY event_count DESC
    LIMIT 3 )
 SELECT ts.sport_name, ts.event_count AS medal_events_count FROM top_sports ts UNION ALL SELECT 'TOTAL' AS sport_name, SUM(ts.event_count) AS medal_events_count FROM  top_sports ts`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := RetrieveScript(test.input); got != test.want {
				t.Errorf("got %s, \nwant %s", got, test.want)
			}
		})

	}
}
