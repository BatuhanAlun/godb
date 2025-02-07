package godb

func CreateDB(dbName, path string) Database {
	return Database{Name: dbName, Path: path, Tables: make([]*Table, 0)}
}

func (d *Database) AddTable(tableName *Table) {
	tableName.DBPath = d.Path
	d.Tables = append(d.Tables, tableName)

}
