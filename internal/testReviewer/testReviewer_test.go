package testReviewer

import "testing"

func TestRetrieveScript(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"with ;", "select * from query;", "select * from query;"},
		{"without ;", "select * from query", "select * from query"},
		{"double quotemarks", `"select * from query"`, "select * from query"},
		{"single quotemarks", `'select * from query'`, "select * from query"},
		{"parentheses", `(select * from query)`, "select * from query"},
		{"square brackets", `[select * from query]`, "select * from query"},
		{"nested query", `select * from (select id from query) s limit 1;`, "select * from (select id from query) s limit 1;"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := RetrieveScript(test.input); got != test.want {
				t.Errorf("got %s, want %s", got, test.want)
			}
		})

	}
}
