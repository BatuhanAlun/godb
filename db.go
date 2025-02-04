package godb

func CreateDB(dbName, path string) Database {
	return Database{Name: dbName, Path: path, Tables: make([]Table, 0)}
}
