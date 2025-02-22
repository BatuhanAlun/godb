package godb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

func (t *Table) Get() []map[string]interface{} {
	var tempMap []map[string]interface{}
	for _, v := range t.Rows {
		tempMap = append(tempMap, v.Data)
	}
	return tempMap
}

func (t *Table) Delete(findCol string, findVal interface{}) error {
	colIndex := -1
	for index, v := range t.Columns {
		if v.Name == findCol {
			colIndex = index
		}
	}
	if colIndex == -1 {
		return fmt.Errorf("wrong column name")
	}

	newRows := t.Rows[:0]
	for _, row := range t.Rows {
		if row.Data[findCol] != findVal {
			newRows = append(newRows, row)
		}
	}

	t.Rows = newRows
	return nil
}

func LoadDatabaseFromFile(dbName string) (*Database, error) {
	dbPath := "/"

	// Read the directory
	files, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database directory: %v", err)
	}

	// Initialize the database
	db := &Database{
		Name:   dbName,
		Path:   dbPath,
		Tables: []*Table{},
	}

	// Iterate over each file in the directory
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue // Skip non-JSON files
		}

		// Table name is the file name without the .json extension
		tableName := strings.TrimSuffix(file.Name(), ".json")
		filePath := filepath.Join(dbPath, dbName, file.Name())

		// Read the table data
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read table file '%s': %v", filePath, err)
		}

		// Deserialize the rows from the JSON file
		var rows []*Row
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("failed to decode JSON for table '%s': %v", tableName, err)
		}

		// Create a table and load its rows
		table := &Table{
			Name:    tableName,
			DBPath:  dbPath,
			Rows:    rows,
			Columns: []Column{}, // Here you could also dynamically load columns from each row if necessary
		}

		// Append the table to the database
		db.Tables = append(db.Tables, table)

		// Debugging: Print to verify that tables are being added
		fmt.Printf("Table '%s' has been added to the database.\n", table.Name)
	}

	// Debugging: Print the total number of tables loaded
	fmt.Printf("Total Tables Loaded: %d\n", len(db.Tables))

	return db, nil
}
