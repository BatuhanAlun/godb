package godb

import (
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

func (d *Database) SaveDatabaseToFile() {

}
