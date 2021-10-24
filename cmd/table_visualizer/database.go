package main

type Database interface {
	Connect(user string, password string, host string, database string)
	Close()
	Select(query string, args ...interface{}) (Table, error)
	GetTables() ([]Table, error)
	GetTable(tableName string) (Table, error)
	GetFks(tableName string, colName string) ([]Fk, error)
}
