package godb

type Database struct {
	Name   string
	Tables []*Table
	Path   string
}

type Table struct {
	Name   string
	DBPath string
	Rows   []any
}
