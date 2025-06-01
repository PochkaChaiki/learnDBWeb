package scriptManager

type ScriptRunResult struct {
	Error   string   `json:"error"`
	Columns []string `json:"columns"`
	Data    [][]any  `json:"data"`
}

type ScriptOptions struct {
	Sql        string
	DBMS       string
	SchemaName string
	TaskName   string
}
