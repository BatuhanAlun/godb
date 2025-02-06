package godb

func CreateTable(tableName string) *Table {
	return &Table{Name: tableName}
}

func CreateRow(data map[string]interface{}) Row {
	return Row{Data: data}
}

func CreateColumn(cName, cType string, cMode ...string) Column {
	mode := ""
	if len(cMode) > 0 && cMode[0] == "PK" {
		mode = "PK"
	}
	return Column{Name: cName, Type: cType, Mode: mode}
}

func (t *Table) AddRow(row *Row) {
	t.Rows = append(t.Rows, row)
}
