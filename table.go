package godb

import (
	"errors"
	"fmt"
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
	} else {
		mode = ""
	}
	if cType != "int" && cType != "bool" && cType != "string" && cType != "UUID" {
		log.Fatal(errors.New("type needs to be int, bool, string or UUID "))
	}
	return Column{Name: cName, Type: cType, Mode: mode}
}

func (t *Table) AddData(mapKey []string, mapVal []interface{}) error {
	if len(mapKey) != len(mapVal) {
		return fmt.Errorf("number of keys and values do not match")
	}

	tempMap := make(map[string]interface{})

	for i, key := range mapKey {
		tempMap[key] = mapVal[i]
	}

	if err := ValidateData(t.Columns, tempMap); err != nil {
		return err
	}

	tempRow := &Row{Data: tempMap}
	t.Rows = append(t.Rows, tempRow)

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

func (t *Table) Get() []map[string]interface{} {
	var tempMap []map[string]interface{}
	for _, v := range t.Rows {
		tempMap = append(tempMap, v.Data)
	}
	return tempMap
}

func (r *Row) UpdateRow(newRow Row) {
	for key, value := range newRow.Data {
		r.Data[key] = value
	}
}
