package godb

type Database struct {
	Name   string
	Tables []*Table
	Path   string
}

type Table struct {
	Name    string
	DBPath  string
	Rows    []*Row
	Columns []Column
}

type Row struct {
	Data map[string]interface{}
}

type Column struct {
	Name string
	Type string
	Mode string
}
