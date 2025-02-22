package godb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ValidateData(columns []Column, row map[string]interface{}) error {
	for _, col := range columns {
		value, exists := row[col.Name]
		if !exists {
			return fmt.Errorf("missing value for column: %s", col.Name)
		}

		switch col.Type {
		case "int":
			_, ok := value.(int)
			_, bok := value.([]uint8)
			if !ok && !bok {
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
	fmt.Println("Saving database...")

	// Ensure db.Path is not empty and directory exists
	if db.Path == "" {
		return fmt.Errorf("db.Path is not set")
	}

	if err := os.MkdirAll(db.Path, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create database folder: %v", err)
	}

	if len(db.Tables) == 0 {
		return fmt.Errorf("no tables found in the database")
	}

	for _, table := range db.Tables {
		if len(table.Rows) == 0 {
			continue // Skip empty tables or handle as needed
		}

		// Validate rows in the table
		for _, row := range table.Rows {
			if err := ValidateData(table.Columns, row.Data); err != nil {
				return fmt.Errorf("validation failed in table '%s': %v", table.Name, err)
			}
		}

		// Marshal table rows to JSON
		tableData, err := json.MarshalIndent(table.Rows, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to encode table '%s' to JSON: %v", table.Name, err)
		}

		// Construct the correct file path
		tablename := table.Name + ".json"
		filePath := filepath.Join(db.Path, tablename)
		fmt.Println("Saving to path:", filePath)

		// Write table data to the file
		err = os.WriteFile(filePath, tableData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write table '%s' to file: %v", table.Name, err)
		}
	}
	return nil
}
