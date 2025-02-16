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

func (t *Table) Update(findCol string, findVal interface{}, updateCol string, updateVal interface{}) error {
	rowIndex := -1

	colExists := false
	for _, col := range t.Columns {
		if col.Name == findCol {
			colExists = true
			break
		}
	}
	if !colExists {
		return fmt.Errorf("column %s does not exist", findCol)
	}

	for index, row := range t.Rows {
		if dataVal, exists := row.Data[findCol]; exists {
			if reflect.TypeOf(dataVal) != reflect.TypeOf(findVal) {
				return fmt.Errorf("type mismatch: expected %T but got %T for %s", dataVal, findVal, findCol)
			}
			if dataVal == findVal {
				rowIndex = index
				break
			}
		}
	}

	if rowIndex == -1 {
		return fmt.Errorf("no matching row found for %s = %v", findCol, findVal)
	}

	colExists = false
	for _, col := range t.Columns {
		if col.Name == updateCol {
			colExists = true
			break
		}
	}
	if !colExists {
		return fmt.Errorf("update column %s does not exist", updateCol)
	}

	existingVal, exists := t.Rows[rowIndex].Data[updateCol]
	if !exists {
		return fmt.Errorf("column %s does not exist in the row", updateCol)
	}

	if reflect.TypeOf(existingVal) != reflect.TypeOf(updateVal) {
		return fmt.Errorf("type mismatch: expected %T but got %T for %s", existingVal, updateVal, updateCol)
	}

	t.Rows[rowIndex].Data[updateCol] = updateVal
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
