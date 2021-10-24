package main

type Column struct {
	Name     string
	Type     string
	Key      string
	Default  string
	Extra    string
	Nullable bool
	Fks      []Fk
}

type Table struct {
	Name    string
	Columns []Column
	Data    [][]string
}

type Fk struct {
	RefCol string
	Table  string
}
