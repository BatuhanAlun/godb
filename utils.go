package godb

import (
	"encoding/json"
	"fmt"
	"os"
)

func (row *Row) ValidateData(columns []Column) error {
	for _, col := range columns {
		value, exists := row.Data[col.Name]
		if !exists {
			return fmt.Errorf("missing value for column: %s", col.Name)
		}

		switch col.Type {
		case "int":
			if _, ok := value.(int); !ok {
				return fmt.Errorf("Column %s expects an int, but got %T", col.Name, value)
			}
		case "string":
			if _, ok := value.(string); !ok {
				return fmt.Errorf("Column %s expects an string, but got %T", col.Name, value)
			}
		case "bool":
			if _, ok := value.(string); !ok {
				return fmt.Errorf("Column %s expects an bool, but got %T", col.Name, value)
			}
		}
	}
	return nil
}

func (d *Database) CreateFiles() error {

	dbFolderPath := d.Path + d.Name
	err := os.Mkdir(dbFolderPath, 0755)
	if err != nil {
		return fmt.Errorf("DB folder cannot be created")
	}

	for _, v := range d.Tables {
		file, err := os.Create(dbFolderPath + "/" + v.Name + ".json")
		if err != nil {
			return fmt.Errorf("error on Creating Table json file ")
		}
		defer file.Close()
	}
	return nil
}

func (db *Database) SaveDatabaseToFile() error {
	if err := os.MkdirAll(db.Name, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create database folder: %v", err)
	}

	for _, table := range db.Tables {
		for _, row := range table.Rows {
			if err := row.ValidateData(table.Columns); err != nil {
				return fmt.Errorf("validation failed in table '%s': %v", table.Name, err)
			}
		}

		tableData, err := json.MarshalIndent(table.Rows, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to encode table '%s' to JSON: %v", table.Name, err)
		}

		filePath := db.Path + db.Name + "/" + table.Name + ".json"

		err = os.WriteFile(filePath, tableData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write table '%s' to file: %v", table.Name, err)
		}
	}

	return nil
}
