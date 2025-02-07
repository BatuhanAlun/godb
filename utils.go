package godb

import (
	"fmt"
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
