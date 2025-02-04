package godb

func CreateTable(tableName string) *Table {
	return &Table{Name: tableName}
}

func CreateRow(rowName, rowMode, rowType string) *Row {
	return &Row{Name: rowName, Mode: rowMode, Type: rowType}
}

func (t *Table) AddRow(row *Row) {
	t.Rows = append(t.Rows, row)
}
