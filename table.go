package godb

import (
	"errors"
	"log"
	"reflect"
	"slices"
)

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

	if cType != "int" && cType != "bool" && cType != "string" {
		log.Fatal(errors.New("type needs to be int, bool or string"))
	}

	return Column{Name: cName, Type: cType, Mode: mode}
}

func (t *Table) AddRow(row *Row) error {
	if err := row.ValidateData(t.Columns); err != nil {
		return err
	}
	t.Rows = append(t.Rows, row)
	return nil
}

func (t *Table) AddColumn(col Column) {
	for _, v := range t.Columns {
		if v.Name == col.Name {
			log.Fatal(errors.New("Column Names Needs to Be Unique"))
		}
	}
	t.Columns = append(t.Columns, col)
}

func (t *Table) DeleteRow(row Row) {
	for index, v := range t.Rows {
		if reflect.DeepEqual(v.Data, row.Data) {
			t.Rows = slices.Delete(t.Rows, index, index+1)
		}
	}
}
